package postgres

import (
	"context"
	"fmt"
	"math"
	"regexp"
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
	FetchExact(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error)
	FetchFuzzy(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error)
	Find(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error)
	FindSimilar(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error)
}

func (r *SearchRepo) FetchExact(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	ids := make([]int, len(req.Items))
	names := make([]string, len(req.Items))
	qtys := make([]float64, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.Id
		names[i] = item.Name
		qtys[i] = item.Quantity
	}

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                id as req_item_id,
                name as req_name,
                qty as req_qty
            FROM UNNEST($1::text[], $2::numeric[], $3::int[]) AS t(name, qty, id)
        )
        SELECT 
            o.id::text, o.year, o.customer, o.consumer, o.date,
            r.req_item_id, p.id::text as matched_item_id, 
			CASE WHEN p.search=r.req_name THEN p.search ELSE p.normalized_notes END as p_search,
            r.req_qty, p.quantity as db_qty,
            1.0 as similarity -- Для точного поиска всегда 1.0
        FROM req r
        JOIN %s p ON p.search = r.req_name OR p.normalized_notes = r.req_name
        JOIN %s o ON o.id = p.order_id`,
		Tables.Positions, Tables.Orders,
	)

	rows, err := r.db.Query(ctx, query, names, qtys, ids)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.RawMatch
	for rows.Next() {
		var m models.RawMatch
		err := rows.Scan(
			&m.OrderId, &m.YearInt, &m.Customer, &m.Consumer, &m.Date,
			&m.ReqId, &m.PosId, &m.PSearch, &m.ReqQty, &m.DbQty, &m.Similarity,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &m)
	}

	return results, nil
}

func (r *SearchRepo) FetchFuzzy(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	ids := make([]int, len(req.Items))
	names := make([]string, len(req.Items))
	qtys := make([]float64, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.Id
		names[i] = item.Name
		qtys[i] = item.Quantity
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, _ = tx.Exec(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.3;`)

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                id as req_id,
                name as req_name,
                qty as req_qty,
                regexp_split_to_array(name, '[^a-zA-Zа-яА-Я0-9.-]+') as req_tokens
            FROM UNNEST($1::text[], $2::numeric[], $3::int[]) AS t(name, qty, id)
        )
        SELECT 
            o.id::text, o.year, o.customer, o.consumer, o.date,
            r.req_id::text,
            p.id::text as pos_id,
            p.search,
            r.req_tokens,
            r.req_qty,
            p.quantity as db_qty,
            word_similarity(r.req_name, p.search) as sml
        FROM req r
        JOIN %s p ON p.search %% r.req_name 
        JOIN %s o ON o.id = p.order_id;`,
		Tables.Positions, Tables.Orders,
	)

	rows, err := tx.Query(ctx, query, names, qtys, ids)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*models.RawMatch
	for rows.Next() {
		var m models.RawMatch
		// Порядок Scan должен строго соответствовать порядку колонок в SELECT
		err := rows.Scan(
			&m.OrderId,
			&m.YearInt,
			&m.Customer,
			&m.Consumer,
			&m.Date,
			&m.ReqId,
			&m.PosId,
			&m.PSearch,
			&m.ReqTokens,
			&m.ReqQty,
			&m.DbQty,
			&m.Similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results, &m)
	}

	return results, nil
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
                o.id as order_id,
                o.customer,
                o.consumer,
                o.year,
                r.req_item_id,
                p.id as matched_item_id,
                0.8 + (0.2 * (LEAST(p.quantity, r.req_qty) / GREATEST(p.quantity, r.req_qty))) as item_score
            FROM req r
            JOIN %s p ON p.search = r.req_name
            JOIN %s o ON o.id = p.order_id
        ),
        order_stats AS (
            SELECT 
                order_id, customer, consumer, year,
                array_agg(matched_item_id) AS matched_item_ids,
                SUM(item_score) as total_score_points
            FROM matches
            GROUP BY order_id, customer, consumer, year
        )
        SELECT 
            order_id, year, customer, consumer,
            matched_item_ids,
            cardinality(matched_item_ids) as matched_names_count,
            cardinality($1) as total_req_count,
            ROUND(
                ((total_score_points / cardinality($1)) * 100)::numeric, 
                2
            ) AS score
        FROM order_stats
        ORDER BY score DESC, year DESC;`,
		Tables.Positions, Tables.Orders,
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
			// &match.PositionIds, // TODO: убрать если я буду использовать этот запрос
			&match.MatchedPos,
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

	_, _ = tx.Exec(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.3;`)

	query := fmt.Sprintf(`
        WITH req AS (
            SELECT 
                idx as req_id,
                name as req_name,
                qty as req_qty,
                regexp_split_to_array(name, '[^a-zA-Zа-яА-Я0-9.-]+') as req_tokens
            FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
        )
        SELECT 
            o.id::text, o.year, o.customer, o.consumer,
            r.req_id::text,
            p.id::text as pos_id,
            p.search,
            r.req_tokens,
            r.req_qty,
            p.quantity as db_qty,
            word_similarity(r.req_name, p.search) as sml
        FROM req r
        JOIN %s p ON p.search %% r.req_name 
        JOIN %s o ON o.id = p.order_id;`, // Убрали WHERE по количеству
		Tables.Positions, Tables.Orders,
	)

	rows, err := tx.Query(ctx, query, names, qtys)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	type matchInfo struct {
		posID     string
		itemScore float64 // Совокупный балл (текст + количество)
	}

	type orderInfo struct {
		Id       string
		Year     int
		Customer string
		Consumer string
		Matches  map[string]matchInfo
	}

	orderMap := make(map[string]*orderInfo)

	for rows.Next() {
		var year int
		var oID, reqID, posID, pSearch, customer, consumer string
		var reqTokens []string
		var reqQty, dbQty, sml float64

		if err := rows.Scan(&oID, &year, &customer, &consumer, &reqID, &posID, &pSearch, &reqTokens, &reqQty, &dbQty, &sml); err != nil {
			return nil, err
		}

		// --- ВАЛИДАЦИЯ ТОКЕНОВ ---
		isValid := true
		pSearchLower := strings.ToLower(pSearch)
		dbTokens := r.makeTokenMap(pSearchLower)

		matchedSpecificCount := 0
		totalSpecificCount := 0

		for _, token := range reqTokens {
			// Считаем токен "специфичным", если в нем есть цифры или он длинный
			isSpecific := r.containsDigits(token) || len(token) > 5

			if isSpecific {
				totalSpecificCount++
				if _, exists := dbTokens[token]; exists {
					matchedSpecificCount++
				} else {
					// Если специфичный токен (например, часть артикула) не найден - это мусор
					isValid = false
					break
				}
			} else if len(token) >= 2 {
				// Для коротких слов (типа "Г" или "В") требуем точного совпадения,
				// но не бракуем весь результат сразу, если это общее слово.
				// Но если это тип СНП (А, В, Г), он критичен.
				if r.isCriticalType(token) {
					if _, exists := dbTokens[token]; !exists {
						isValid = false
						break
					}
				}
			}
		}

		if !isValid || (totalSpecificCount > 0 && matchedSpecificCount == 0) {
			continue
		}

		// --- РАСЧЕТ ВЕСА ПОЗИЦИИ ---
		// Базовый вес = текстовое сходство (0.3 - 1.0)
		// Множитель количества: если совпадает идеально = 1.0, если разница в 10 раз = ~0.8
		// Формула: 80% за имя + 20% за близость количества
		ratio := math.Min(reqQty, dbQty) / math.Max(reqQty, dbQty)
		qtyFactor := 0.8 + (0.2 * ratio)

		currentItemScore := sml * qtyFactor

		if _, ok := orderMap[oID]; !ok {
			orderMap[oID] = &orderInfo{
				Id: oID, Year: year, Customer: customer, Consumer: consumer,
				Matches: make(map[string]matchInfo),
			}
		}

		// Сохраняем только лучшее совпадение для конкретной строки запроса
		if prev, exists := orderMap[oID].Matches[reqID]; !exists || currentItemScore > prev.itemScore {
			orderMap[oID].Matches[reqID] = matchInfo{posID: posID, itemScore: currentItemScore}
		}
	}

	totalReqCount := float64(len(req.Items))
	results := make([]*models.OrderMatchResult, 0, len(orderMap))

	for _, info := range orderMap {
		var sumScore float64
		posIDs := make([]string, 0, len(info.Matches))

		for _, m := range info.Matches {
			sumScore += m.itemScore
			posIDs = append(posIDs, m.posID)
		}

		// Итоговый Score: среднее качество совпадения по всем позициям запроса
		finalScore := math.Round((sumScore/totalReqCount*100)*100) / 100

		// Выводим заказ, если есть ХОТЯ БЫ ОДНО совпадение (finalScore > 0)
		if finalScore > 0 {
			results = append(results, &models.OrderMatchResult{
				OrderId:  info.Id,
				Year:     info.Year,
				Customer: info.Customer,
				Consumer: info.Consumer,
				// PositionIds: posIDs, // TODO: убрать если я буду использовать этот запрос
				MatchedPos: len(info.Matches),
				TotalCount: int(totalReqCount),
				Score:      finalScore,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Year > results[j].Year
	})

	return results, tx.Commit(ctx)
}

// Создаем карту токенов из строки БД для моментальной проверки
func (r *SearchRepo) makeTokenMap(searchStr string) map[string]struct{} {
	words := strings.Fields(strings.ToLower(searchStr))
	m := make(map[string]struct{}, len(words))
	for _, w := range words {
		m[w] = struct{}{}
	}
	return m
}
func (r *SearchRepo) containsDigits(s string) bool {
	return strings.ContainsAny(s, "0123456789")
}

var reSymbols = regexp.MustCompile(`[а-яА-Яa-zA-Z]`)

// Критичные токены - это одиночные буквы, которые часто значат тип (А, В, Г, П)
func (r *SearchRepo) isCriticalType(s string) bool {
	if len(s) != 1 {
		return false
	}
	// Проверяем, что это буква, а не просто мусор
	return reSymbols.MatchString(s)
}

// func (r *SearchRepo) FindSimilar(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
// 	names := make([]string, len(req.Items))
// 	qtys := make([]float64, len(req.Items))

// 	for i, item := range req.Items {
// 		names[i] = item.Name
// 		qtys[i] = item.Quantity
// 	}

// 	tx, err := r.db.Begin(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback(ctx)

// 	// Настраиваем чувствительность триграмм
// 	_, err = tx.Exec(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.6;`)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := fmt.Sprintf(`
// 		WITH req AS (
// 			SELECT
// 				idx as req_item_id,
// 				name as req_name,
// 				qty as req_qty,
// 				-- Извлекаем все числа из запроса в массив
// 				regexp_split_to_array(regexp_replace(name, '\D', ' ', 'g'), '\s+') as req_numbers
// 			FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
// 		),
// 		matches AS (
// 			SELECT
// 				o.id, o.customer, o.consumer, o.year,
// 				r.req_item_id,
// 				p.id as matched_item_id,
// 				-- Считаем финальный скор
// 				word_similarity(r.req_name, p.search) as sim_score
// 			FROM req r
// 			-- 1. Быстрый поиск по индексу триграмм
// 			JOIN %s p ON p.search %% r.req_name
// 			JOIN %s o ON o.id = p.order_id
// 			WHERE
// 				-- 2. Фильтр по количеству
// 				p.quantity BETWEEN r.req_qty * 0.7 AND r.req_qty * 1.3
// 				-- 3. ЦИФРОВОЙ КОНТРОЛЛЕР:
// 				-- Проверяем, что каждое число >= 2 знаков из запроса есть в строке базы
// 				AND NOT EXISTS (
// 					SELECT FROM unnest(r.req_numbers) AS n
// 					WHERE length(n) >= 2 AND p.search NOT LIKE '%%' || n || '%%'
// 				)
// 		),
// 		order_stats AS (
// 			SELECT
// 				id, customer, consumer, year,
// 				array_agg(matched_item_id) AS matched_item_ids,
// 				COUNT(DISTINCT req_item_id) AS matched_req_count
// 			FROM matches
// 			GROUP BY id, customer, consumer, year
// 		)
// 		SELECT
// 			id, year, customer, consumer,
// 			matched_item_ids,
// 			matched_req_count,
// 			cardinality($1) as total_req_count,
// 			ROUND((matched_req_count::numeric / cardinality($1)) * 100, 2) AS score
// 		FROM order_stats
// 		WHERE (matched_req_count::numeric / cardinality($1)) >= 0.70
// 		ORDER BY year DESC, score DESC;`,
// 		Tables.Positions, Tables.Orders,
// 	)

// 	rows, err := tx.Query(ctx, query, names, qtys)
// 	if err != nil {
// 		return nil, fmt.Errorf("query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	var results []*models.OrderMatchResult
// 	for rows.Next() {
// 		match := &models.OrderMatchResult{}
// 		if err := rows.Scan(
// 			&match.OrderId, &match.Year, &match.Customer, &match.Consumer,
// 			&match.PositionIds, &match.MatchedCount, &match.TotalCount, &match.Score,
// 		); err != nil {
// 			return nil, err
// 		}
// 		results = append(results, match)
// 	}

// 	return results, tx.Commit(ctx)
// }

// func (r *SearchRepo) FindSimilar2(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
// 	names := make([]string, len(req.Items))
// 	qtys := make([]float64, len(req.Items))

// 	for i, item := range req.Items {
// 		names[i] = item.Name
// 		qtys[i] = item.Quantity
// 	}

// 	tx, err := r.db.Begin(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback(ctx)

// 	// Порог чувствительности делаем ниже, чтобы индекс находил больше кандидатов
// 	_, err = tx.Exec(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.3;`)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := fmt.Sprintf(`
// 		WITH req AS (
// 			SELECT
// 				idx as req_item_id,
// 				name as req_name,
// 				qty as req_qty,
// 				regexp_split_to_array(regexp_replace(name, '\D', ' ', 'g'), '\s+') as req_numbers
// 			FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
// 		),
// 		raw_matches AS (
// 			SELECT
// 				o.id as order_id, o.customer, o.consumer, o.year,
// 				r.req_item_id,
// 				p.id as p_id,
// 				-- Считаем скор для конкретной пары (запрос <-> позиция в базе)
// 				(word_similarity(r.req_name, p.search) * (0.8 + 0.2 * CASE WHEN r.req_qty = 0 THEN 1.0 ELSE (LEAST(p.quantity, r.req_qty) / GREATEST(p.quantity, r.req_qty)) END) * CASE WHEN NOT EXISTS (
// 					SELECT FROM unnest(r.req_numbers) AS n
// 					WHERE length(n) >= 2 AND p.search NOT LIKE '%%' || n || '%%'
// 				) THEN 1.0 ELSE 0.5 END) as item_score
// 			FROM req r
// 			JOIN %s p ON p.search %% r.req_name
// 			JOIN %s o ON o.id = p.order_id
// 		),
// 		ranked_matches AS (
// 			SELECT
// 				*,
// 				-- Важный момент: если для одной позиции запроса нашлось 5 в базе,
// 				-- мы нумеруем их по убыванию качества
// 				ROW_NUMBER() OVER(PARTITION BY order_id, req_item_id ORDER BY item_score DESC) as rank
// 			FROM raw_matches
// 		)
// 		SELECT
// 			order_id, year, customer, consumer,
// 			-- Собираем только те ID, которые заняли 1-е место по качеству для своей позиции
// 			array_agg(p_id) FILTER (WHERE rank = 1) as matched_item_ids,
// 			COUNT(req_item_id) FILTER (WHERE rank = 1) as matched_count,
// 			cardinality($1) as total_count,
// 			ROUND(((SUM(item_score) FILTER (WHERE rank = 1) / cardinality($1)) * 100)::numeric, 2) AS score
// 		FROM ranked_matches
// 		GROUP BY order_id, year, customer, consumer
// 		ORDER BY score DESC, year DESC
// 		LIMIT 100`,
// 		Tables.Positions, Tables.Orders,
// 	)

// 	rows, err := tx.Query(ctx, query, names, qtys)
// 	if err != nil {
// 		return nil, fmt.Errorf("query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	var results []*models.OrderMatchResult
// 	for rows.Next() {
// 		match := &models.OrderMatchResult{}
// 		if err := rows.Scan(
// 			&match.OrderId, &match.Year, &match.Customer, &match.Consumer,
// 			&match.PositionIds, &match.MatchedCount, &match.TotalCount, &match.Score,
// 		); err != nil {
// 			return nil, err
// 		}
// 		results = append(results, match)
// 	}

// 	return results, tx.Commit(ctx)
// }

// func (r *SearchRepo) FindSimilar2(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
// 	names := make([]string, len(req.Items))
// 	qtys := make([]float64, len(req.Items))

// 	for i, item := range req.Items {
// 		names[i] = item.Name
// 		qtys[i] = item.Quantity
// 	}

// 	tx, err := r.db.Begin(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback(ctx)

// 	// Настраиваем чувствительность триграмм
// 	_, err = tx.Exec(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.5;`)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := fmt.Sprintf(`
//         WITH req AS (
//             SELECT
//                 idx as req_id,
//                 name as req_name,
//                 qty as req_qty,
//                 regexp_split_to_array(regexp_replace(name, '\D', ' ', 'g'), '\s+') as req_digits
//             FROM UNNEST($1::text[], $2::numeric[]) WITH ORDINALITY AS t(name, qty, idx)
//         )
//         SELECT
//             o.id, o.year, o.customer, o.consumer,
//             r.req_id,
//             p.id as pos_id,
//             p.search,
//             r.req_digits
//         FROM req r
//         JOIN %s p ON p.search %% r.req_name
//         JOIN %s o ON o.id = p.order_id
//         WHERE p.quantity BETWEEN r.req_qty * 0.7 AND r.req_qty * 1.3;`,
// 		Tables.Positions, Tables.Orders,
// 	)

// 	rows, err := tx.Query(ctx, query, names, qtys)
// 	if err != nil {
// 		return nil, fmt.Errorf("query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	type orderInfo struct {
// 		Id       string
// 		Year     int
// 		Customer string
// 		Consumer string
// 		Matches  map[string]string // req_id -> pos_id
// 	}

// 	orderMap := make(map[string]*orderInfo)

// 	for rows.Next() {
// 		var year int
// 		var oID, reqID, posID string
// 		var customer, consumer, pSearch string
// 		var reqDigits []string

// 		if err := rows.Scan(&oID, &year, &customer, &consumer, &reqID, &posID, &pSearch, &reqDigits); err != nil {
// 			return nil, err
// 		}

// 		// --- ЦИФРОВОЙ КОНТРОЛЛЕР (Фильтрация в Go) ---
// 		// Проверяем, что все числа длиной >= 2 из запроса есть в названии позиции
// 		isValid := true
// 		for _, digit := range reqDigits {
// 			if len(digit) >= 2 && !strings.Contains(pSearch, digit) {
// 				isValid = false
// 				break
// 			}
// 		}

// 		if !isValid {
// 			continue
// 		}

// 		// Группируем результаты по заказам
// 		if _, ok := orderMap[oID]; !ok {
// 			orderMap[oID] = &orderInfo{
// 				Id:       oID,
// 				Year:     year,
// 				Customer: customer,
// 				Consumer: consumer,
// 				Matches:  make(map[string]string),
// 			}
// 		}
// 		orderMap[oID].Matches[reqID] = posID
// 	}

// 	// Формируем финальный список результатов
// 	totalReqCount := len(req.Items)
// 	results := make([]*models.OrderMatchResult, 0)

// 	for _, info := range orderMap {
// 		matchedCount := len(info.Matches)
// 		score := math.Round((float64(matchedCount)/float64(totalReqCount)*100)*100) / 100

// 		// Условие: минимум 70% совпадений позиций в одном заказе
// 		if score >= 70.0 {
// 			posIDs := make([]string, 0, len(info.Matches))
// 			for _, pid := range info.Matches {
// 				posIDs = append(posIDs, pid)
// 			}

// 			results = append(results, &models.OrderMatchResult{
// 				OrderId:      info.Id,
// 				Year:         info.Year,
// 				Customer:     info.Customer,
// 				Consumer:     info.Consumer,
// 				PositionIds:  posIDs,
// 				MatchedCount: matchedCount,
// 				TotalCount:   totalReqCount,
// 				Score:        score,
// 			})
// 		}
// 	}

// 	// Сортировка: Сначала по проценту совпадения, затем по году
// 	sort.Slice(results, func(i, j int) bool {
// 		if results[i].Score != results[j].Score {
// 			return results[i].Score > results[j].Score
// 		}
// 		return results[i].Year > results[j].Year
// 	})

// 	return results, tx.Commit(ctx)
// }
