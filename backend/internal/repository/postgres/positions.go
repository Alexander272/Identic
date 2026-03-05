package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PositionRepo struct {
	db *pgxpool.Pool
	Transaction
}

func NewPositionRepo(db *pgxpool.Pool, tr Transaction) *PositionRepo {
	return &PositionRepo{
		db:          db,
		Transaction: tr,
	}
}

type Positions interface {
	Create(ctx context.Context, tx Tx, dto []*models.PositionDTO) error
	Update(ctx context.Context, tx Tx, dto []*models.PositionDTO) error
	Delete(ctx context.Context, tx Tx, dto []*models.DeletePositionDTO) error
}

func (r *PositionRepo) Create(ctx context.Context, tx Tx, dto []*models.PositionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	ids := make([]string, len(dto))
	orderIds := make([]string, len(dto))
	names := make([]string, len(dto))
	quantities := make([]float64, len(dto))
	notes := make([]string, len(dto))

	for i, v := range dto {
		if v.Id == "" {
			v.Id = uuid.NewString()
		}

		ids[i] = v.Id
		orderIds[i] = v.OrderId
		names[i] = v.Name
		quantities[i] = v.Quantity
		notes[i] = v.Notes
	}

	query := fmt.Sprintf(`INSERT INTO %s (id, order_id, name, quantity, notes)
		SELECT unnest($1::uuid[]), unnest($2::uuid[]), unnest($3::text[]), unnest($4::real[]), unnest($5::text[])`,
		PositionsTable,
	)

	if _, err := r.getExec(tx).Exec(ctx, query, ids, orderIds, names, quantities, notes); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionRepo) Update(ctx context.Context, tx Tx, dto []*models.PositionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	ids := make([]string, len(dto))
	names := make([]string, len(dto))
	quantities := make([]float64, len(dto))
	notes := make([]string, len(dto))

	for i, v := range dto {
		ids[i] = v.Id
		names[i] = v.Name
		quantities[i] = v.Quantity
		notes[i] = v.Notes
	}

	query := fmt.Sprintf(`UPDATE %s SET name = $2, quantity = $3, notes = $4
		SET name = s.name, quantity = s.quantity, notes = s.notes 
		FROM (
			SELECT UNNEST($1::uuid[]) as id, 
				   UNNEST($2::uuid[]) as name, 
				   UNNEST($3::real[]) as quantity,
				   UNNEST($4::text[]) as notes
		) AS s 
		WHERE t.id = s.id`,
		PositionsTable,
	)

	if _, err := r.getExec(tx).Exec(ctx, query, ids, names, quantities, notes); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionRepo) Delete(ctx context.Context, tx Tx, dto []*models.DeletePositionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	ids := make([]string, len(dto))
	for i, v := range dto {
		ids[i] = v.Id
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ANY($1::uuid[])`, PositionsTable)

	if _, err := r.getExec(tx).Exec(ctx, query, ids); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
