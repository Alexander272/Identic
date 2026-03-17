package postgres

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
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
	"customer": {},
	"consumer": {},
	"manager":  {},
	"year":     {},
	"date":     {},
	"notes":    {},
}

type Orders interface {
	GetById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error)
	GetByYear(ctx context.Context, req *models.GetOrderByYearDTO) ([]*models.Order, error)
	GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error)
	Create(ctx context.Context, tx Tx, dto *models.OrderDTO) error
	CreateSeveral(ctx context.Context, tx Tx, dto []*models.OrderDTO) error
	Update(ctx context.Context, tx Tx, dto *models.OrderDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteOrderDTO) error
}

func (r *OrderRepo) GetById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error) {
	query := fmt.Sprintf(`SELECT id, customer, consumer, manager, bill, date, notes FROM %s WHERE id = $1`,
		OrdersTable,
	)
	order := &models.Order{}

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&order.Id, &order.Customer, &order.Consumer, &order.Manager, &order.Bill, &order.Date, &order.Notes,
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
	query := fmt.Sprintf(`SELECT o.id, o.customer, o.consumer, o.manager, o.bill, o.date, o.notes, COUNT(p.id) AS position_count
        FROM %s AS o
        LEFT JOIN %s AS p ON p.order_id = o.id
        WHERE o.year = $1
        GROUP BY o.id ORDER BY o.date DESC, manager, customer, consumer`,
		OrdersTable, PositionsTable,
	)
	var data []*models.Order

	rows, err := r.db.Query(ctx, query, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.Order{}
		if err := rows.Scan(&tmp.Id, &tmp.Customer, &tmp.Consumer, &tmp.Manager, &tmp.Bill, &tmp.Date, &tmp.Notes, &tmp.PositionCount); err != nil {
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
	snake := reCamelCase.ReplaceAllString(req.Field, "${1}_${2}")
	req.Field = strings.ToLower(snake)

	if _, exist := allowedFields[req.Field]; !exist {
		return nil, models.ErrFieldNotAllowed
	}
	quotedField := pgx.Identifier{req.Field}.Sanitize()

	query := fmt.Sprintf(`SELECT COALESCE(array_agg(DISTINCT %s::text), '{}'::text[]) 
		FROM %s WHERE %s::text!='' AND %s IS NOT NULL`,
		quotedField, OrdersTable, quotedField, quotedField,
	)
	var data []string

	err := r.db.QueryRow(ctx, query).Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	slices.SortFunc(data, func(a, b string) int {
		if req.Sort == "DESC" {
			return strings.Compare(b, a) // Убывание
		}
		return strings.Compare(a, b) // Возрастание
	})

	return data, nil
}

func (r *OrderRepo) Create(ctx context.Context, tx Tx, dto *models.OrderDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, customer, consumer, manager, bill, date, year, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		OrdersTable,
	)
	if dto.Id == "" {
		dto.Id = uuid.NewString()
	}

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.Id, dto.Customer, dto.Consumer, dto.Manager, dto.Bill, dto.Date, dto.Date.Year(), dto.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) CreateSeveral(ctx context.Context, tx Tx, dto []*models.OrderDTO) error {
	if len(dto) == 0 {
		return nil
	}

	// ids := make([]string, len(dto))
	// customers := make([]string, len(dto))
	// consumers := make([]string, len(dto))
	// managers := make([]string, len(dto))
	// bills := make([]string, len(dto))
	// dates := make([]time.Time, len(dto))
	// notes := make([]string, len(dto))

	// for i, v := range dto {
	// 	ids[i] = v.Id
	// 	customers[i] = v.Customer
	// 	consumers[i] = v.Consumer
	// 	managers[i] = v.Manager
	// 	bills[i] = v.Bill
	// 	dates[i] = v.Date
	// 	notes[i] = v.Notes
	// }

	// query := fmt.Sprintf(`INSERT INTO %s (id, customer, consumer, manager, bill, date, notes)
	// 	SELECT unnest($1::uuid[]), unnest($2::text[]), unnest($3::text[]), unnest($4::text[]),
	// 		unnest($5::text[]), unnest($6::timestamp with time zone[]), unnest($7::text[])`,
	// 	OrdersTable,
	// )

	// if _, err := r.getExec(tx).Exec(ctx, query, ids, customers, consumers, managers, bills, dates, notes); err != nil {
	// 	return fmt.Errorf("failed to execute query. error: %w", err)
	// }
	// return nil

	rows := make([][]interface{}, len(dto))
	for i, v := range dto {
		rows[i] = []interface{}{
			v.Id,
			v.Customer,
			v.Consumer,
			v.Manager,
			v.Bill,
			v.Date,
			v.Date.Year(),
			v.Notes,
		}
	}

	columns := []string{"id", "customer", "consumer", "manager", "bill", "date", "year", "notes"}
	_, err := r.getExec(tx).CopyFrom(
		ctx,
		pgx.Identifier{OrdersTable},
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Update(ctx context.Context, tx Tx, dto *models.OrderDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET customer = $2, consumer = $3, manager = $4, bill = $5, date = $6, year = $7, notes = $8 
		WHERE id = $1`,
		OrdersTable,
	)

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.Id, dto.Customer, dto.Consumer, dto.Manager, dto.Bill, dto.Date, dto.Date.Year(), dto.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteOrderDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, OrdersTable)
	_, err := r.getExec(tx).Exec(ctx, query, dto.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
