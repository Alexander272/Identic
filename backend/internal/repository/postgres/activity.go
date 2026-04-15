package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
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
	GetByEntity(ctx context.Context, req *models.GetActivityLogsDTO) ([]*models.ActivityLog, error)
	GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*models.ActivityLog, error)
}

func (r *ActivityRepo) Create(ctx context.Context, tx Tx, dto *models.CreateActivityLogDTO) error {
	query := `INSERT INTO activity_logs (action, changed_by, changed_by_name, entity_type, entity_id, parent_id, old_values, new_values)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	var oldJSON, newJSON interface{}
	if dto.OldValues != nil {
		oldJSON = dto.OldValues
	}
	if dto.NewValues != nil {
		newJSON = dto.NewValues
	}

	_, err := r.getExec(tx).Exec(ctx, query,
		dto.Action, dto.ChangedBy, dto.ChangedByName,
		dto.EntityType, dto.EntityID, dto.ParentID,
		oldJSON, newJSON,
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

func (r *ActivityRepo) GetByEntity(ctx context.Context, req *models.GetActivityLogsDTO) ([]*models.ActivityLog, error) {
	var query string
	var args []interface{}

	if req.ParentID != nil {
		query = `SELECT id, action, changed_by, changed_by_name, entity_type, entity_id, parent_id, old_values, new_values, created_at
			FROM activity_logs WHERE parent_id = $1 ORDER BY created_at DESC`
		args = append(args, req.ParentID)
	} else if req.EntityID != "" {
		query = `SELECT id, action, changed_by, changed_by_name, entity_type, entity_id, parent_id, old_values, new_values, created_at
			FROM activity_logs WHERE entity_id = $1 AND entity_type = $2 ORDER BY created_at DESC`
		args = append(args, req.EntityID, req.EntityType)
	} else {
		return nil, fmt.Errorf("either entity_id or parent_id must be provided")
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity logs: %w", err)
	}
	defer rows.Close()

	logs := make([]*models.ActivityLog, 0, 20)
	for rows.Next() {
		log := &models.ActivityLog{}
		var oldBytes, newBytes []byte

		err := rows.Scan(&log.ID, &log.Action, &log.ChangedBy, &log.ChangedByName,
			&log.EntityType, &log.EntityID, &log.ParentID, &oldBytes, &newBytes, &log.CreatedAt)
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

func (r *ActivityRepo) GetByOrder(ctx context.Context, orderID uuid.UUID) ([]*models.ActivityLog, error) {
	query := `SELECT id, action, changed_by, changed_by_name, entity_type, entity_id, parent_id, old_values, new_values, created_at
		FROM activity_logs WHERE entity_id = $1 OR parent_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity logs: %w", err)
	}
	defer rows.Close()

	logs := make([]*models.ActivityLog, 0, 20)
	for rows.Next() {
		log := &models.ActivityLog{}
		var oldBytes, newBytes []byte

		err := rows.Scan(&log.ID, &log.Action, &log.ChangedBy, &log.ChangedByName,
			&log.EntityType, &log.EntityID, &log.ParentID, &oldBytes, &newBytes, &log.CreatedAt)
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
