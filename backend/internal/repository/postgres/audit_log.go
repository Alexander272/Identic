package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type auditRepo struct {
	db *pgxpool.Pool
	Transaction
}

func NewAuditRepo(db *pgxpool.Pool, tr Transaction) *auditRepo {
	return &auditRepo{
		db:          db,
		Transaction: tr,
	}
}

type AuditLogs interface {
	Get(ctx context.Context, req *models.GetAuditLogsDTO) ([]*models.AuditLog, error)
	Create(ctx context.Context, tx Tx, dto *models.AuditLogDTO) error
	CreateSeveral(ctx context.Context, tx Tx, dto []*models.AuditLogDTO) error
}

func (r *auditRepo) Get(ctx context.Context, req *models.GetAuditLogsDTO) ([]*models.AuditLog, error) {
	//TODO наверное мне какие-нибудь фильтры или ограничения понадобятся
	query := fmt.Sprintf(`SELECT id, changed_by, changed_by_name, action, entity_type, entity_id, 
		old_values, new_values, created_at 
		FROM %s`,
		Tables.AuditLogs,
	)
	var data []*models.AuditLog

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.AuditLog{}
		if err := rows.Scan(
			&tmp.ID, &tmp.ChangedBy, &tmp.ChangedByName, &tmp.Action, &tmp.EntityType, &tmp.EntityID,
			&tmp.OldValues, &tmp.NewValues, &tmp.CreatedAt,
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

func (r *auditRepo) Create(ctx context.Context, tx Tx, dto *models.AuditLogDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, changed_by, changed_by_name, action, entity_type, entity_id, old_values, new_values) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		Tables.AuditLogs,
	)
	dto.ID = uuid.New()

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.ID, dto.ChangedBy, dto.ChangedByName, dto.Action, dto.EntityType, dto.EntityID, dto.OldValues, dto.NewValues,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *auditRepo) CreateSeveral(ctx context.Context, tx Tx, dto []*models.AuditLogDTO) error {
	if len(dto) == 0 {
		return nil
	}

	rows := make([][]interface{}, len(dto))

	for i, v := range dto {
		if v.ID == uuid.Nil {
			v.ID = uuid.New()
		}

		rows[i] = []interface{}{
			v.ID,
			v.ChangedBy,
			v.ChangedByName,
			v.Action,
			v.EntityType,
			v.EntityID,
			v.OldValues,
			v.NewValues,
		}
	}

	columns := []string{"id", "changed_by", "changed_by_name", "action", "entity_type", "entity_id", "old_values", "new_values"}
	_, err := r.getExec(tx).CopyFrom(
		ctx,
		pgx.Identifier{Tables.AuditLogs},
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
