package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
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

type Orders interface {
	Create(ctx context.Context, tx Tx, dto *models.OrderDTO) error
	CreateSeveral(ctx context.Context, tx Tx, dto []*models.OrderDTO) error
	Update(ctx context.Context, tx Tx, dto *models.OrderDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteOrderDTO) error
}

func (r *OrderRepo) Create(ctx context.Context, tx Tx, dto *models.OrderDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, customer, consumer, manager, bill, date, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		OrdersTable,
	)
	if dto.Id == "" {
		dto.Id = uuid.NewString()
	}

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.Id, dto.Customer, dto.Consumer, dto.Manager, dto.Bill, dto.Date, dto.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to create order. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) CreateSeveral(ctx context.Context, tx Tx, dto []*models.OrderDTO) error {
	if len(dto) == 0 {
		return nil
	}

	ids := make([]string, len(dto))
	customers := make([]string, len(dto))
	consumers := make([]string, len(dto))
	managers := make([]string, len(dto))
	bills := make([]string, len(dto))
	dates := make([]time.Time, len(dto))
	notes := make([]string, len(dto))

	for i, v := range dto {
		ids[i] = v.Id
		customers[i] = v.Customer
		consumers[i] = v.Consumer
		managers[i] = v.Manager
		bills[i] = v.Bill
		dates[i] = v.Date
		notes[i] = v.Notes
	}

	query := fmt.Sprintf(`INSERT INTO %s (id, customer, consumer, manager, bill, date, notes)
		SELECT unnest($1::uuid[]), unnest($2::text[]), unnest($3::text[]), unnest($4::text[]), 
			unnest($5::text[]), unnest($6::timestamp with time zone[]), unnest($7::text[])`,
		OrdersTable,
	)

	if _, err := r.getExec(tx).Exec(ctx, query, ids, customers, consumers, managers, bills, dates, notes); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Update(ctx context.Context, tx Tx, dto *models.OrderDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET customer = $2, consumer = $3, manager = $4, bill = $5, date = $6, notes = $7 
		WHERE id = $1`,
		OrdersTable,
	)

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.Id, dto.Customer, dto.Consumer, dto.Manager, dto.Bill, dto.Date, dto.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to update order. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteOrderDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, OrdersTable)
	_, err := r.getExec(tx).Exec(ctx, query, dto.Id)
	if err != nil {
		return fmt.Errorf("failed to delete order. error: %w", err)
	}
	return nil
}
