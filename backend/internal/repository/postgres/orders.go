package postgres

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	db *pgxpool.Pool
	Transaction
}

func NewOrderRepo(db *pgxpool.Pool, tr Transaction) *OrderRepo {
	return &OrderRepo{
		db:          db,
		Transaction: tr,
	}
}

var reCamelCase = regexp.MustCompile("([a-z0-9])([0-9A-Z])")
var allowedFields = map[string]struct{}{
	"customer":      {},
	"consumer":      {},
	"client":        {},
	"manager":       {},
	"year":          {},
	"date":          {},
	"notes":         {},
	"is_bargaining": {},
	"is_budget":     {},
}

type Orders interface {
	Get(ctx context.Context, req *models.OrderFilterDTO) ([]*models.Order, error)
	GetById(ctx context.Context, tx Tx, req *models.GetOrderByIdDTO) (*models.Order, error)
	GetByYear(ctx context.Context, req *models.GetOrderByYearDTO) ([]*models.Order, error)
	GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error)
	GetFlatData(ctx context.Context, req *models.GetFlatOrderDTO) (*models.FlatOrderRes, error)
	IsExist(ctx context.Context, tx Tx, dto *models.OrderDTO) (bool, error)
	IsExistByPos(ctx context.Context, tx Tx, dto *models.OrderDTO) (bool, error)
	Create(ctx context.Context, tx Tx, dto *models.OrderDTO) error
	CreateSeveral(ctx context.Context, tx Tx, dto []*models.OrderDTO) error
	Update(ctx context.Context, tx Tx, dto *models.OrderDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteOrderDTO) error
}

func (r *OrderRepo) Get(ctx context.Context, req *models.OrderFilterDTO) ([]*models.Order, error) {
	var allowedFields = map[string][]string{
		"client":       {"customer", "consumer"},
		"manager":      {"manager"},
		"date":         {"date"},
		"isBargaining": {"is_bargaining"},
		"isBudget":     {"is_budget"},
	}

	baseQuery := fmt.Sprintf(`SELECT o.id, o.customer, o.consumer, o.manager, o.is_bargaining, o.is_budget, 
		o.bill, o.date, o.notes, COUNT(p.id) AS position_count
        FROM %s AS o
        LEFT JOIN %s AS p ON p.order_id = o.id`,
		Tables.Orders, Tables.Positions,
	)
	qb := NewQueryBuilder(baseQuery)

	for _, filter := range req.Filters {
		cols := allowedFields[filter.Field]

		for _, val := range filter.Values {
			if len(cols) > 1 {
				types := slices.Repeat([]string{val.CompareType}, len(cols))
				values := slices.Repeat([]string{val.Value}, len(cols))
				qb.AddCompositeFilter(cols, types, values)
			} else if len(cols) == 1 {
				qb.AddFilter(cols[0], val.CompareType, val.Value)
			}
		}
	}

	qb.SetGroupBy("o.id")

	sortFields := []SortField{
		{Field: "o.date", Desc: true},
		{Field: "o.created_at", Desc: true},
		{Field: "manager", Desc: false},
		{Field: "customer", Desc: false},
		{Field: "consumer", Desc: false},
	}
	sortFields = append(sortFields, SortField{Field: "o.id", Desc: false})
	qb.SetMultiSort(sortFields)

	query, args := qb.Build()
	var data []*models.Order

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.Order{}
		if err := rows.Scan(
			&tmp.Id, &tmp.Customer, &tmp.Consumer, &tmp.Manager, &tmp.IsBargaining, &tmp.IsBudget,
			&tmp.Bill, &tmp.Date, &tmp.Notes, &tmp.PositionCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row. error: %w", err)
		}
		data = append(data, tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return data, nil
}

func (r *OrderRepo) GetById(ctx context.Context, tx Tx, req *models.GetOrderByIdDTO) (*models.Order, error) {
	query := fmt.Sprintf(`SELECT id, customer, consumer, manager, is_bargaining, is_budget, bill, date, notes, created_at FROM %s WHERE id = $1`,
		Tables.Orders,
	)
	order := &models.Order{}

	err := r.getExec(tx).QueryRow(ctx, query, req.Id).Scan(
		&order.Id, &order.Customer, &order.Consumer, &order.Manager, &order.IsBargaining, &order.IsBudget,
		&order.Bill, &order.Date, &order.Notes, &order.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return order, nil
}

func (r *OrderRepo) GetByYear(ctx context.Context, req *models.GetOrderByYearDTO) ([]*models.Order, error) {
	query := fmt.Sprintf(`SELECT o.id, o.customer, o.consumer, o.manager, o.is_bargaining, o.is_budget, 
		o.bill, o.date, o.notes, COUNT(p.id) AS position_count
        FROM %s AS o
        LEFT JOIN %s AS p ON p.order_id = o.id
        WHERE o.year = $1
        GROUP BY o.id ORDER BY o.date DESC, o.created_at DESC, manager, customer, consumer`,
		Tables.Orders, Tables.Positions,
	)
	var data []*models.Order

	rows, err := r.db.Query(ctx, query, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.Order{}
		if err := rows.Scan(
			&tmp.Id, &tmp.Customer, &tmp.Consumer, &tmp.Manager, &tmp.IsBargaining, &tmp.IsBudget,
			&tmp.Bill, &tmp.Date, &tmp.Notes, &tmp.PositionCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row. error: %w", err)
		}
		data = append(data, tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return data, nil
}

func (r *OrderRepo) GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error) {
	var validFields []string

	rawFields := []string{req.Field}
	if req.Field == "client" {
		rawFields = []string{"customer", "consumer"}
	}

	for _, f := range rawFields {
		f = strings.TrimSpace(f)
		snake := reCamelCase.ReplaceAllString(f, "${1}_${2}")
		field := strings.ToLower(snake)

		if _, exist := allowedFields[field]; !exist {
			return nil, models.ErrFieldNotAllowed
		}
		validFields = append(validFields, pgx.Identifier{field}.Sanitize())
	}

	// 2. Строим подзапрос, который объединяет значения из всех выбранных колонок в одну
	// Используем LATERAL unnest для эффективного превращения колонок в строки
	columnList := strings.Join(validFields, ", ")

	// cross join lateral позволяет для каждой строки таблицы
	// превратить список колонок в набор строк
	query := fmt.Sprintf(`
		SELECT COALESCE(array_agg(DISTINCT val::text), '{}'::text[])
		FROM %s,
		LATERAL (SELECT unnest(ARRAY[%s])) AS t(val)
		WHERE val::text != '' AND val IS NOT NULL`,
		Tables.Orders, columnList,
	)

	var data []string
	err := r.db.QueryRow(ctx, query).Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	slices.SortFunc(data, func(a, b string) int {
		if req.Sort == "DESC" {
			return strings.Compare(strings.ToLower(b), strings.ToLower(a))
		}
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	return data, nil
}

func (r *OrderRepo) GetFlatData(ctx context.Context, req *models.GetFlatOrderDTO) (*models.FlatOrderRes, error) {
	baseQuery := fmt.Sprintf(`SELECT p.id, customer, consumer, manager, is_bargaining, is_budget, bill, date, o.notes, 
		row_number, name, quantity, p.notes AS pos_notes
		FROM %s AS o
		JOIN %s AS p ON p.order_id = o.id`,
		Tables.Orders, Tables.Positions,
	)
	qb := NewQueryBuilder(baseQuery)

	searchFields := []string{}

	if req.Search != nil {
		for _, f := range req.Search.Fields {
			searchFields = append(searchFields, pgx.Identifier{f}.Sanitize())
		}
		qb.AddMultiSearch(searchFields, req.Search.Value) // поиск
	}

	sortFields := []SortField{
		{Field: "date", Desc: true},
		{Field: "customer", Desc: false},
		{Field: "consumer", Desc: false},
		{Field: "row_number", Desc: false},
	}

	if req.Sort != nil {
		sortFields = []SortField{
			{Field: pgx.Identifier{req.Sort.Field}.Sanitize(), Desc: req.Sort.Type == "DESC"},
		}
	}
	sortFields = append(sortFields, SortField{Field: "p.id", Desc: false})
	qb.SetMultiSort(sortFields) // сортировка

	// parts := strings.Split(req.Cursor, "|")
	// if len(parts) == 2 {
	// 	cursorDate := parts[0]
	// 	cursorID := parts[1]
	// 	qb.SetCompositeCursor(sortField, cursorDate, cursorID, sortDesc)
	// }
	// parts := strings.Split(req.Cursor, "|")
	// if len(parts) >= 2 {
	// 	cursorDate := parts[0]
	// 	cursorOrderID := parts[1]

	// 	fields := []string{"date", "p.id"}
	// 	values := []interface{}{cursorDate, cursorOrderID}
	// 	descSlice := []bool{true, false} // должно совпадать с сортировкой!

	// 	qb.SetMultiCompositeCursor(fields, values, descSlice)
	// }

	if req.Cursor != "" {
		cursorState, err := DecodeCursor(req.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}

		// Проверяем, что курсор соответствует текущей сортировке
		// (опционально, но полезно для отладки)
		// ...

		// Парсим значения с учётом типов
		values, err := cursorState.ParseCursorValues()
		if err != nil {
			return nil, fmt.Errorf("failed to parse cursor: %w", err)
		}

		// Извлекаем поля и направления из курсора или берём из sortConfig
		fields := make([]string, 0, len(cursorState.Types))
		desc := cursorState.Desc
		if len(desc) == 0 {
			// fallback: берём из sortConfig
			for _, sf := range sortFields {
				desc = append(desc, sf.Desc)
			}
		}
		// Поля должны совпадать с теми, что в сортировке
		// Здесь можно добавить валидацию

		// Для простоты: берём поля из sortConfig
		for _, sf := range sortFields {
			fields = append(fields, sf.Field)
		}

		qb.SetMultiCompositeCursor(fields, values, desc)
	}

	limit := 0
	if req.Page != nil {
		limit = req.Page.Limit
		qb.SetLimit(limit + 1)
	}

	query, args := qb.Build()

	var data []*models.FlatOrder

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.FlatOrder{}
		if err := rows.Scan(
			&tmp.Id, &tmp.Customer, &tmp.Consumer, &tmp.Manager, &tmp.IsBargaining, &tmp.IsBudget, &tmp.Bill, &tmp.Date, &tmp.Notes,
			&tmp.RowNumber, &tmp.Name, &tmp.Quantity, &tmp.PositionNotes,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row. error: %w", err)
		}
		data = append(data, tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	hasMore := false
	if limit > 0 && len(data) > limit {
		hasMore = true
		data = data[:limit] // отрезаем лишнюю запись
	}

	// var lastDate time.Time
	// var lastID string
	// if len(data) > 0 {
	// 	lastDate = data[len(data)-1].Date
	// 	lastID = data[len(data)-1].Id
	// }

	var nextCursor string
	if hasMore && len(data) > 0 {
		last := data[len(data)-1]

		// Формируем значения в порядке полей сортировки
		rowValues := make([]interface{}, 0, len(sortFields))
		fieldTypes := make([]string, 0, len(sortFields))
		desc := make([]bool, 0, len(sortFields))

		for _, sf := range sortFields {
			val, typ, ok := last.CursorValue(sf.Field)
			if !ok {
				return nil, fmt.Errorf("unknown cursor field: %s", sf.Field)
			}

			rowValues = append(rowValues, val)
			fieldTypes = append(fieldTypes, typ)
			desc = append(desc, sf.Desc)
		}

		nextCursor, err = BuildCursorFromRow(rowValues, fieldTypes, desc)
		if err != nil {
			return nil, fmt.Errorf("failed to build cursor: %w", err)
		}
	}

	res := &models.FlatOrderRes{
		Orders:  data,
		Cursor:  nextCursor,
		HasMore: hasMore,
	}
	return res, nil
}

func (r *OrderRepo) IsExist(ctx context.Context, tx Tx, dto *models.OrderDTO) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS (
			SELECT 1 
			FROM %s AS o
			INNER JOIN %s AS p ON o.id = p.order_id
			WHERE o.customer = $1 
			AND o.consumer = $2 
			AND o.notes = $3 
			AND date >= $4 AND date < $4 + interval '1 day'
			GROUP BY o.id
			HAVING COUNT(p.id) = $5
		)`,
		Tables.Orders, Tables.Positions,
	)
	var exists bool

	err := r.db.QueryRow(ctx, query, dto.Customer, dto.Consumer, dto.Notes, dto.Date, len(dto.Positions)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return exists, nil
}
func (r *OrderRepo) IsExistByPos(ctx context.Context, tx Tx, dto *models.OrderDTO) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS (
			SELECT 1 FROM %s 
			WHERE content_hash = $1 AND created_at > NOW() - INTERVAL '14 days'
		)`,
		Tables.Orders,
	)

	var exists bool

	err := r.db.QueryRow(ctx, query, dto.Hash).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return exists, nil
}

func (r *OrderRepo) Create(ctx context.Context, tx Tx, dto *models.OrderDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, customer, consumer, manager, is_bargaining, is_budget, bill, date, year, notes, content_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		Tables.Orders,
	)
	if dto.Id == "" {
		dto.Id = uuid.NewString()
	}
	dto.Year = dto.Date.Year()

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.Id, dto.Customer, dto.Consumer, dto.Manager, dto.IsBargaining, dto.IsBudget, dto.Bill, dto.Date, dto.Year, dto.Notes, dto.Hash,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	payload, _ := json.Marshal(ws_hub.WSMessage{
		Action: "ORDER_INSERTED",
		Data: models.Order{
			Id:            dto.Id,
			Customer:      dto.Customer,
			Consumer:      dto.Consumer,
			Manager:       dto.Manager,
			IsBargaining:  dto.IsBargaining,
			IsBudget:      dto.IsBudget,
			Bill:          dto.Bill,
			Date:          dto.Date,
			Year:          dto.Year,
			Notes:         dto.Notes,
			PositionCount: len(dto.Positions),
		},
	})
	_, err = r.getExec(tx).Exec(ctx, "SELECT pg_notify('order_updates', $1)", string(payload))
	if err != nil {
		return fmt.Errorf("failed to execute notify query. error: %w", err)
	}

	return nil
}

func (r *OrderRepo) CreateSeveral(ctx context.Context, tx Tx, dto []*models.OrderDTO) error {
	if len(dto) == 0 {
		return nil
	}

	rows := make([][]interface{}, len(dto))
	yearsMap := make(map[int]struct{})
	for i, v := range dto {
		yearsMap[v.Date.Year()] = struct{}{}
		rows[i] = []interface{}{
			v.Id,
			v.Customer,
			v.Consumer,
			v.Manager,
			v.IsBargaining,
			v.IsBudget,
			v.Bill,
			v.Date,
			v.Date.Year(),
			v.Notes,
		}
	}

	var affectedYears []int
	for y := range yearsMap {
		affectedYears = append(affectedYears, y)
	}

	columns := []string{"id", "customer", "consumer", "manager", "is_bargaining", "is_budget", "bill", "date", "year", "notes"}
	_, err := r.getExec(tx).CopyFrom(
		ctx,
		pgx.Identifier{Tables.Orders},
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	// 3. Отправка уведомления через pg_notify
	payload, _ := json.Marshal(ws_hub.WSMessage{
		Action: "ORDER_BULK_INSERTED",
		Data: map[string]interface{}{
			"years": affectedYears,
		},
	})

	// Выполняем NOTIFY. Важно делать это в той же транзакции (tx),
	// чтобы уведомление ушло только если транзакция закоммитится.
	notifyQuery := fmt.Sprintf("SELECT pg_notify('order_updates', %s)", quoteLiteral(string(payload)))
	if _, err = r.getExec(tx).Exec(ctx, notifyQuery); err != nil {
		return fmt.Errorf("failed to execute notify query. error: %w", err)
	}

	return nil
}

func quoteLiteral(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (r *OrderRepo) Update(ctx context.Context, tx Tx, dto *models.OrderDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET customer = $2, consumer = $3, manager = $4, is_bargaining = $5, is_budget = $6, 
		bill = $7, date = $8, year = $9, notes = $10, content_hash=$11
		WHERE id = $1`,
		Tables.Orders,
	)
	dto.Year = dto.Date.Year()

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.Id, dto.Customer, dto.Consumer, dto.Manager, dto.IsBargaining, dto.IsBudget, dto.Bill, dto.Date, dto.Year, dto.Notes, dto.Hash,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	payload, _ := json.Marshal(ws_hub.WSMessage{
		Action: "ORDER_UPDATED",
		Data: models.Order{
			Id:            dto.Id,
			Customer:      dto.Customer,
			Consumer:      dto.Consumer,
			Manager:       dto.Manager,
			IsBargaining:  dto.IsBargaining,
			IsBudget:      dto.IsBudget,
			Bill:          dto.Bill,
			Date:          dto.Date,
			Year:          dto.Year,
			Notes:         dto.Notes,
			PositionCount: len(dto.Positions),
		},
	})

	_, err = r.getExec(tx).Exec(ctx, "SELECT pg_notify('order_updates', $1)", string(payload))
	if err != nil {
		return fmt.Errorf("failed to execute notify query. error: %w", err)
	}

	return nil
}

func (r *OrderRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteOrderDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 RETURNING year`, Tables.Orders)
	var year int

	err := r.getExec(tx).QueryRow(ctx, query, dto.Id).Scan(&year)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // Заказ уже удален или не существовал
		}
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	payload, _ := json.Marshal(ws_hub.WSMessage{
		Action: "ORDER_DELETED",
		Data: map[string]interface{}{
			"id":   dto.Id,
			"year": year,
		},
	})

	_, err = r.getExec(tx).Exec(ctx, "SELECT pg_notify('order_updates', $1)", string(payload))
	if err != nil {
		return fmt.Errorf("failed to execute notify query. error: %w", err)
	}

	return nil
}
