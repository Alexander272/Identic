package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	GetByOrder(ctx context.Context, req *models.GetPositionsByOrderIdDTO) ([]*models.Position, error)
	GetByIds(ctx context.Context, req *models.GetPositionsByIds) ([]*models.Position, error)
	Create(ctx context.Context, tx Tx, dto []*models.PositionDTO) error
	Update(ctx context.Context, tx Tx, dto []*models.PositionDTO) error
	Delete(ctx context.Context, tx Tx, dto []*models.DeletePositionDTO) error
}

func (r *PositionRepo) GetByOrder(ctx context.Context, req *models.GetPositionsByOrderIdDTO) ([]*models.Position, error) {
	query := fmt.Sprintf(`SELECT id, order_id, row_number, name, quantity, notes FROM %s WHERE order_id = $1 ORDER BY row_number ASC`,
		PositionsTable,
	)

	var positions []*models.Position
	rows, err := r.db.Query(ctx, query, req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.Position{}
		if err := rows.Scan(&tmp.Id, &tmp.OrderId, &tmp.RowNumber, &tmp.Name, &tmp.Quantity, &tmp.Notes); err != nil {
			return nil, fmt.Errorf("failed to scan row. error: %w", err)
		}
		positions = append(positions, tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	// for i  := range positions {
	// 	positions[i].Quantity = math.Round(positions[i].Quantity * 100) / 100
	// }

	return positions, nil
}

func (r *PositionRepo) GetByIds(ctx context.Context, req *models.GetPositionsByIds) ([]*models.Position, error) {
	query := fmt.Sprintf(`SELECT id, order_id, row_number, name, quantity, notes FROM %s WHERE id = ANY($1::uuid[]) ORDER BY row_number ASC`,
		PositionsTable,
	)

	positions := make([]*models.Position, 0, len(req.Ids))
	rows, err := r.db.Query(ctx, query, req.Ids)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.Position{}
		if err := rows.Scan(&tmp.Id, &tmp.OrderId, &tmp.RowNumber, &tmp.Name, &tmp.Quantity, &tmp.Notes); err != nil {
			return nil, fmt.Errorf("failed to scan row. error: %w", err)
		}
		positions = append(positions, tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return positions, nil
}

func (r *PositionRepo) Create(ctx context.Context, tx Tx, dto []*models.PositionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	// ids := make([]string, len(dto))
	// orderIds := make([]string, len(dto))
	// names := make([]string, len(dto))
	// searches := make([]string, len(dto))
	// quantities := make([]float64, len(dto))
	// notes := make([]string, len(dto))

	// for i, v := range dto {
	// 	if v.Id == "" {
	// 		v.Id = uuid.NewString()
	// 	}

	// 	ids[i] = v.Id
	// 	orderIds[i] = v.OrderId
	// 	names[i] = v.Name
	// 	searches[i] = v.Search
	// 	quantities[i] = v.Quantity
	// 	notes[i] = v.Notes
	// }

	// query := fmt.Sprintf(`INSERT INTO %s (id, order_id, name, search, quantity, notes)
	// 	SELECT unnest($1::uuid[]), unnest($2::uuid[]), unnest($3::text[]), unnest($4::text[]), unnest($5::real[]), unnest($6::text[])`,
	// 	PositionsTable,
	// )

	// if _, err := r.getExec(tx).Exec(ctx, query, ids, orderIds, names, searches, quantities, notes); err != nil {
	// 	return fmt.Errorf("failed to execute query. error: %w", err)
	// }
	// return nil

	rows := make([][]interface{}, len(dto))

	for i, v := range dto {
		if v.Id == "" {
			v.Id = uuid.NewString()
		}

		rows[i] = []interface{}{
			v.Id,
			v.OrderId,
			v.RowNumber,
			v.Name,
			v.Search,
			v.Quantity,
			v.Notes,
		}
	}

	columns := []string{"id", "order_id", "row_number", "name", "search", "quantity", "notes"}
	_, err := r.getExec(tx).CopyFrom(
		ctx,
		pgx.Identifier{PositionsTable},
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
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
	search := make([]string, len(dto))
	quantities := make([]float32, len(dto))
	notes := make([]string, len(dto))

	for i, v := range dto {
		ids[i] = v.Id
		names[i] = v.Name
		quantities[i] = v.Quantity
		notes[i] = v.Notes
	}

	query := fmt.Sprintf(`UPDATE %s AS t SET name = s.name, search = s.search, quantity = s.quantity, notes = s.notes 
		FROM (
			SELECT UNNEST($1::uuid[]) as id, 
				   UNNEST($2::text[]) as name, 
				   UNNEST($3::text[]) as search,
				   UNNEST($4::real[]) as quantity,
				   UNNEST($5::text[]) as notes
		) AS s 
		WHERE t.id = s.id`,
		PositionsTable,
	)

	if _, err := r.getExec(tx).Exec(ctx, query, ids, names, search, quantities, notes); err != nil {
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
