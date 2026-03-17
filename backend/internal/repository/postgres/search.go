package postgres

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchRepo struct {
	db *pgxpool.Pool
}

func NewSearchRepo(db *pgxpool.Pool) *SearchRepo {
	return &SearchRepo{
		db: db,
	}
}

type Search interface {
	Find(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error)
	FindSimilar(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error)
	FindSimilar2(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error)
}

func (r *SearchRepo) Find(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
	names := make([]string, len(req.Items))
	qtys := make([]float64, len(req.Items))

	for i, item := range req.Items {
		names[i] = item.Name
		qtys[i] = item.Quantity
	}

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                idx as req_item_id,
                name as req_name,
                qty as req_qty
            FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
        ),
        matches AS (
            SELECT 
                o.id,
                o.customer,
                o.consumer,
                o.year,
                r.req_item_id,
				p.id as matched_item_id,
				CASE 
					WHEN p.quantity = r.req_qty THEN 1 
					ELSE 0 
				END as qty_match_flag
            FROM req r
            JOIN %s p ON p.search = r.req_name
            JOIN %s o ON o.id = p.order_id
            WHERE 
                -- Допуск по количеству оставляем (разброс 30%%)
                p.quantity BETWEEN r.req_qty * 0.7 AND r.req_qty * 1.3
        ),
        order_stats AS (
            SELECT 
                id, customer, consumer, year,
				array_agg(matched_item_id) AS matched_item_ids,
                COUNT(DISTINCT req_item_id) AS matched_req_count,
				COUNT(DISTINCT CASE WHEN qty_match_flag = 1 THEN req_item_id END) AS total_req_count
            FROM matches
            GROUP BY id, customer, consumer, year
        )
        SELECT 
            id, year, customer, consumer,
			matched_item_ids,
            matched_req_count,
            -- cardinality($1) as total_req_count,
            total_req_count,
            ROUND((matched_req_count::numeric / cardinality($1)) * 100, 2) AS score
        FROM order_stats
        WHERE (matched_req_count::numeric / cardinality($1)) >= 0.70
        ORDER BY year DESC, score DESC;`,
		PositionsTable, OrdersTable,
	)

	rows, err := r.db.Query(ctx, query, names, qtys)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.OrderMatchResult

	for rows.Next() {
		match := &models.OrderMatchResult{}
		err := rows.Scan(
			&match.OrderId,
			&match.Year,
			&match.Customer,
			&match.Consumer,
			&match.PositionIds,
			&match.MatchedCount,
			&match.TotalCount,
			&match.Score,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, match)
	}
	return results, nil
}

func (r *SearchRepo) FindSimilarOld(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
	names := make([]string, len(req.Items))
	qtys := make([]float64, len(req.Items))

	for i, item := range req.Items {
		// Предполагаем, что нормализация уже сделана в коде до этого момента
		names[i] = item.Name
		qtys[i] = item.Quantity
	}

	// query := fmt.Sprintf(`WITH req AS (
	// 		-- Сшиваем два массива в одну таблицу req
	// 		-- ordinality добавит порядковый номер (id), чтобы мы могли различать позиции
	// 		SELECT
	// 			idx as req_item_id,
	// 			name as req_name,
	// 			qty as req_qty
	// 		FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
	// 	),
	// 	matches AS (
	// 		SELECT
	// 			o.id,
	// 			o.customer,
	// 			o.consumer,
	// 			o.year,
	// 			r.req_item_id,
	// 			oi.name AS matched_item_name,
	// 			oi.quantity AS matched_item_qty,
	// 			strict_word_similarity(r.req_name, oi.search) AS sim_score
	// 		FROM req r
	// 		-- Используем оператор %% для задействования GIN индекса
	// 		JOIN %s oi ON oi.search %% r.req_name
	// 		JOIN %s o ON o.id = oi.order_id
	// 		WHERE
	// 			-- Условие по количеству (разброс 30%%)
	// 			oi.quantity BETWEEN r.req_qty * 0.7 AND r.req_qty * 1.3
	// 			AND (
	// 				NOT (r.req_name ~ '\d') OR -- если цифр в запросе нет, пропускаем
	// 				similarity(regexp_replace(r.req_name, '\D', '', 'g'), regexp_replace(oi.search, '\D', '', 'g')) > 0.3
	// 			)
	// 	),
	// 	order_stats AS (
	// 		SELECT
	// 			id,
	// 			customer,
	// 			consumer,
	// 			year,
	// 			COUNT(DISTINCT req_item_id) AS matched_req_count
	// 			-- Собираем только для вывода в Go
	// 			-- jsonb_agg(
	// 			-- 	jsonb_build_object(
	// 			-- 		'name', matched_item_name,
	// 			-- 		'qty', matched_item_qty
	// 			-- 	)
	// 			-- ) AS matched_positions
	// 		FROM matches
	// 		GROUP BY id, customer, consumer, year
	// 	)
	// 	SELECT
	// 		id,
	// 		year,
	// 		customer,
	// 		consumer,
	// 		matched_req_count,
	// 		-- Общее кол-во элементов в массиве $1
	// 		cardinality($1) as total_req_count,
	// 		ROUND((matched_req_count::numeric / cardinality($1)) * 100, 2) AS score
	// 		-- matched_positions
	// 	FROM order_stats
	// 	WHERE (matched_req_count::numeric / cardinality($1)) >= 0.70
	// 	ORDER BY year DESC, score DESC;`,
	// 	PositionsTable, OrdersTable,
	// )
	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                idx as req_item_id,
                name as req_name,
                qty as req_qty,
                -- Извлекаем только цифры из запроса для "цифрового профиля"
                regexp_replace(name, '\D', '', 'g') as req_digits
            FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
        ),
        matches AS (
            SELECT 
                o.id, o.customer, o.consumer, o.year,
                r.req_item_id,
                -- Считаем схожесть всей строки
                strict_word_similarity(r.req_name, oi.search) AS sim_score,
                -- Считаем схожесть только цифровых частей
                similarity(r.req_digits, regexp_replace(oi.search, '\D', '', 'g')) AS digit_sim
            FROM req r
            -- Используем %% для GIN индекса
            JOIN %s oi ON oi.search %% r.req_name
            JOIN %s o ON o.id = oi.order_id
            WHERE 
                oi.quantity BETWEEN r.req_qty * 0.7 AND r.req_qty * 1.3
                AND (
                    -- ЛОГИКА ОТСЕЧЕНИЯ:
                    -- Если в запросе были цифры, то цифровое сходство должно быть высоким
                    -- (например, > 0.4). Это не даст 73 превратиться в 455.
                    r.req_digits = '' OR 
					(strict_word_similarity(r.req_name, oi.search) > 0.5 AND
                    similarity(r.req_digits, regexp_replace(oi.search, '\D', '', 'g')) > 0.4)
                )
        ),
        order_stats AS (
            SELECT 
                id, customer, consumer, year,
                COUNT(DISTINCT req_item_id) AS matched_req_count
            FROM matches
            -- Дополнительно: можно брать только лучшие совпадения для каждой позиции
            WHERE sim_score > 0.5
            GROUP BY id, customer, consumer, year
        )
        SELECT 
            id, year, customer, consumer,
            matched_req_count,
            cardinality($1) as total_req_count,
            ROUND((matched_req_count::numeric / cardinality($1)) * 100, 2) AS score
        FROM order_stats
        WHERE (matched_req_count::numeric / cardinality($1)) >= 0.70
        ORDER BY year DESC, score DESC`,
		PositionsTable, OrdersTable,
	)

	rows, err := r.db.Query(ctx, query, names, qtys)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// var results []*models.Results
	// var currentGroup *models.Results
	var results []*models.OrderMatchResult

	for rows.Next() {
		match := &models.OrderMatchResult{}
		// var posJSON []byte

		err := rows.Scan(
			&match.OrderId,
			&match.Year,
			&match.Customer,
			&match.Consumer,
			&match.MatchedCount,
			&match.TotalCount,
			&match.Score,
			// &posJSON,
		)
		if err != nil {
			return nil, err
		}
		// json.Unmarshal(posJSON, &match.MatchedPositions)

		// Логика группировки по годам (SQL гарантирует сортировку)
		// if currentGroup == nil || currentGroup.Year != year {
		// 	if currentGroup != nil {
		// 		results = append(results, currentGroup)
		// 	}
		// 	currentGroup = &models.Results{Year: year, Orders: []*models.OrderMatchResult{}}
		// }
		// currentGroup.Orders = append(currentGroup.Orders, match)

		results = append(results, match)
	}

	// if currentGroup != nil {
	// 	results = append(results, currentGroup)
	// }

	return results, nil
}

func (r *SearchRepo) FindSimilar(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
	names := make([]string, len(req.Items))
	qtys := make([]float64, len(req.Items))

	for i, item := range req.Items {
		names[i] = item.Name
		qtys[i] = item.Quantity
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Настраиваем чувствительность триграмм
	_, err = tx.Exec(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.6;`)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		WITH req AS (
			SELECT
				idx as req_item_id,
				name as req_name,
				qty as req_qty,
				-- Извлекаем все числа из запроса в массив
				regexp_split_to_array(regexp_replace(name, '\D', ' ', 'g'), '\s+') as req_numbers
			FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
		),
		matches AS (
			SELECT
				o.id, o.customer, o.consumer, o.year,
				r.req_item_id,
				p.id as matched_item_id,
				-- Считаем финальный скор
				word_similarity(r.req_name, p.search) as sim_score
			FROM req r
			-- 1. Быстрый поиск по индексу триграмм
			JOIN %s p ON p.search %% r.req_name
			JOIN %s o ON o.id = p.order_id
			WHERE
				-- 2. Фильтр по количеству
				p.quantity BETWEEN r.req_qty * 0.7 AND r.req_qty * 1.3
				-- 3. ЦИФРОВОЙ КОНТРОЛЛЕР:
				-- Проверяем, что каждое число >= 2 знаков из запроса есть в строке базы
				AND NOT EXISTS (
					SELECT FROM unnest(r.req_numbers) AS n
					WHERE length(n) >= 2 AND p.search NOT LIKE '%%' || n || '%%'
				)
		),
		order_stats AS (
			SELECT
				id, customer, consumer, year,
				array_agg(matched_item_id) AS matched_item_ids,
				COUNT(DISTINCT req_item_id) AS matched_req_count
			FROM matches
			GROUP BY id, customer, consumer, year
		)
		SELECT
			id, year, customer, consumer,
			matched_item_ids,
			matched_req_count,
			cardinality($1) as total_req_count,
			ROUND((matched_req_count::numeric / cardinality($1)) * 100, 2) AS score
		FROM order_stats
		WHERE (matched_req_count::numeric / cardinality($1)) >= 0.70
		ORDER BY year DESC, score DESC;`,
		PositionsTable, OrdersTable,
	)

	rows, err := tx.Query(ctx, query, names, qtys)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.OrderMatchResult
	for rows.Next() {
		match := &models.OrderMatchResult{}
		if err := rows.Scan(
			&match.OrderId, &match.Year, &match.Customer, &match.Consumer,
			&match.PositionIds, &match.MatchedCount, &match.TotalCount, &match.Score,
		); err != nil {
			return nil, err
		}
		results = append(results, match)
	}

	return results, tx.Commit(ctx)
}

func (r *SearchRepo) FindSimilar2(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
	names := make([]string, len(req.Items))
	qtys := make([]float64, len(req.Items))

	for i, item := range req.Items {
		names[i] = item.Name
		qtys[i] = item.Quantity
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Настраиваем чувствительность триграмм
	_, err = tx.Exec(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.5;`)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                idx as req_id,
                name as req_name,
                qty as req_qty,
                regexp_split_to_array(regexp_replace(name, '\D', ' ', 'g'), '\s+') as req_digits
            FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
        )
        SELECT 
            o.id, o.year, o.customer, o.consumer,
            r.req_id,
            p.id as pos_id,
            p.search,
            r.req_digits
        FROM req r
        JOIN %s p ON p.search %% r.req_name 
        JOIN %s o ON o.id = p.order_id
        WHERE p.quantity BETWEEN r.req_qty * 0.7 AND r.req_qty * 1.3;`,
		PositionsTable, OrdersTable,
	)

	rows, err := tx.Query(ctx, query, names, qtys)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	type orderInfo struct {
		Id       string
		Year     int
		Customer string
		Consumer string
		Matches  map[string]string // req_id -> pos_id
	}

	orderMap := make(map[string]*orderInfo)

	for rows.Next() {
		var year int
		var oID, reqID, posID string
		var customer, consumer, pSearch string
		var reqDigits []string

		if err := rows.Scan(&oID, &year, &customer, &consumer, &reqID, &posID, &pSearch, &reqDigits); err != nil {
			return nil, err
		}

		// --- ЦИФРОВОЙ КОНТРОЛЛЕР (Фильтрация в Go) ---
		// Проверяем, что все числа длиной >= 2 из запроса есть в названии позиции
		isValid := true
		for _, digit := range reqDigits {
			if len(digit) >= 2 && !strings.Contains(pSearch, digit) {
				isValid = false
				break
			}
		}

		if !isValid {
			continue
		}

		// Группируем результаты по заказам
		if _, ok := orderMap[oID]; !ok {
			orderMap[oID] = &orderInfo{
				Id:       oID,
				Year:     year,
				Customer: customer,
				Consumer: consumer,
				Matches:  make(map[string]string),
			}
		}
		orderMap[oID].Matches[reqID] = posID
	}

	// Формируем финальный список результатов
	totalReqCount := len(req.Items)
	results := make([]*models.OrderMatchResult, 0)

	for _, info := range orderMap {
		matchedCount := len(info.Matches)
		score := math.Round((float64(matchedCount)/float64(totalReqCount)*100)*100) / 100

		// Условие: минимум 70% совпадений позиций в одном заказе
		if score >= 70.0 {
			posIDs := make([]string, 0, len(info.Matches))
			for _, pid := range info.Matches {
				posIDs = append(posIDs, pid)
			}

			results = append(results, &models.OrderMatchResult{
				OrderId:      info.Id,
				Year:         info.Year,
				Customer:     info.Customer,
				Consumer:     info.Consumer,
				PositionIds:  posIDs,
				MatchedCount: matchedCount,
				TotalCount:   totalReqCount,
				Score:        score,
			})
		}
	}

	// Сортировка: Сначала по проценту совпадения, затем по году
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Year > results[j].Year
	})

	return results, tx.Commit(ctx)
}
