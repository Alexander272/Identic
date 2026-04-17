package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepo struct {
	db *pgxpool.Pool
	Transaction
}

func NewActivityRepo(db *pgxpool.Pool, tr Transaction) *ActivityRepo {
	return &ActivityRepo{
		db:          db,
		Transaction: tr,
	}
}

type Activity interface {
	Create(ctx context.Context, tx Tx, dto *models.CreateActivityLogDTO) error
	CreateBatch(ctx context.Context, tx Tx, dtos []*models.CreateActivityLogDTO) error
	Get(ctx context.Context, req *models.GetAllActivityLogsDTO) ([]*models.ActivityLog, error)
	GetByEntity(ctx context.Context, req *models.GetActivityLogsDTO) ([]*models.ActivityLog, error)
	GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*models.ActivityLog, error)
}

func (r *ActivityRepo) Create(ctx context.Context, tx Tx, dto *models.CreateActivityLogDTO) error {
	query := `INSERT INTO activity_logs (action, changed_by, changed_by_name, entity_type, entity_id, entity, parent_id, old_values, new_values)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	var oldJSON, newJSON interface{}
	if dto.OldValues != nil {
		oldJSON = dto.OldValues
	}
	if dto.NewValues != nil {
		newJSON = dto.NewValues
	}

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.Action, dto.ChangedBy, dto.ChangedByName,
		dto.EntityType, dto.EntityID, dto.Entity,
		dto.ParentID, oldJSON, newJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create activity log: %w", err)
	}
	return nil
}

func (r *ActivityRepo) CreateBatch(ctx context.Context, tx Tx, dtos []*models.CreateActivityLogDTO) error {
	if len(dtos) == 0 {
		return nil
	}

	for _, dto := range dtos {
		if err := r.Create(ctx, tx, dto); err != nil {
			return err
		}
	}
	return nil
}

func (r *ActivityRepo) Get(ctx context.Context, req *models.GetAllActivityLogsDTO) ([]*models.ActivityLog, error) {
	baseQuery := fmt.Sprintf(`SELECT id, action, changed_by, changed_by_name, entity_type, entity_id, entity, 
		parent_id, old_values, new_values, created_at
		FROM %s`,
		Tables.ActivityLogs,
	)

	qb := NewQueryBuilder(baseQuery)
	qb.AddUUIDFilter("actor_id", req.ActorID)
	qb.AddDateRangeFilter("created_at", req.StartDate, req.EndDate)
	qb.SetSort("created_at", true)

	if req.Limit > 0 {
		qb.SetLimit(req.Limit)
	}
	if req.Offset > 0 {
		qb.SetOffset(req.Offset)
	}

	query, args := qb.Build()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity logs: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *ActivityRepo) GetByEntity(ctx context.Context, req *models.GetActivityLogsDTO) ([]*models.ActivityLog, error) {
	baseQuery := fmt.Sprintf(`SELECT id, action, changed_by, changed_by_name, entity_type, entity_id, entity, 
		parent_id, old_values, new_values, created_at
		FROM %s`,
		Tables.ActivityLogs,
	)

	qb := NewQueryBuilder(baseQuery)
	if req.ParentID != nil {
		qb.AddUUIDFilter("parent_id", req.ParentID)
	}
	if req.EntityID != "" {
		qb.AddStringFilter("entity_id", string(req.EntityID))
		qb.AddStringFilter("entity_type", string(req.EntityType))
	}
	qb.SetSort("created_at", true)

	query, args := qb.Build()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity logs: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *ActivityRepo) GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*models.ActivityLog, error) {
	query := `SELECT id, action, changed_by, changed_by_name, entity_type, entity_id, entity, parent_id, old_values, new_values, created_at
		FROM activity_logs WHERE entity_id = $1 OR parent_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity logs: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *ActivityRepo) scanLogs(rows pgx.Rows) ([]*models.ActivityLog, error) {
	logs := make([]*models.ActivityLog, 0, 20)
	for rows.Next() {
		log := &models.ActivityLog{}
		var oldBytes, newBytes []byte

		err := rows.Scan(&log.ID, &log.Action, &log.ChangedBy, &log.ChangedByName,
			&log.EntityType, &log.EntityID, &log.Entity, &log.ParentID, &oldBytes, &newBytes, &log.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity log: %w", err)
		}

		if oldBytes != nil {
			log.OldValues = json.RawMessage(oldBytes)
		}
		if newBytes != nil {
			log.NewValues = json.RawMessage(newBytes)
		}

		logs = append(logs, log)
	}

	return logs, nil
}
