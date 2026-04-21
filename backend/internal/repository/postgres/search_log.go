package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchLogRepo struct {
	db *pgxpool.Pool
}

func NewSearchLogRepo(db *pgxpool.Pool) *SearchLogRepo {
	return &SearchLogRepo{db: db}
}

type SearchLogs interface {
	Create(ctx context.Context, dto *models.CreateSearchLogDTO) error
	Get(ctx context.Context, dto *models.GetSearchLogsDTO) ([]*models.SearchLog, error)
}

func (r *SearchLogRepo) Create(ctx context.Context, dto *models.CreateSearchLogDTO) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (search_id, actor_id, actor_name, search_type, query, duration_ms, results_count, items_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		Tables.SearchLogs,
	)

	queryJSON, err := json.Marshal(dto.Query)
	if err != nil {
		return fmt.Errorf("failed to marshal query: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		dto.SearchId,
		dto.ActorID,
		dto.ActorName,
		dto.SearchType,
		queryJSON,
		dto.DurationMs,
		dto.ResultsCount,
		dto.ItemsCount,
	)
	if err != nil {
		return fmt.Errorf("failed to create search log: %w", err)
	}

	return nil
}

func (r *SearchLogRepo) Get(ctx context.Context, dto *models.GetSearchLogsDTO) ([]*models.SearchLog, error) {
	baseQuery := fmt.Sprintf(`
		SELECT s.id, search_id, actor_id, actor_name, search_type, query, duration_ms, results_count, items_count, s.created_at,
		COALESCE(u.last_name, ''), COALESCE(u.first_name, ''), COALESCE(u.email, '')
		FROM %s s LEFT JOIN %s u ON actor_id::text=u.sso_id`,
		Tables.SearchLogs, Tables.Users,
	)

	qb := NewQueryBuilder(baseQuery)
	qb.AddUUIDFilter("actor_id", dto.ActorID)
	qb.AddDateRangeFilter("s.created_at", dto.StartDate, dto.EndDate)
	qb.SetSort("s.created_at", true)

	if dto.Limit > 0 {
		qb.SetLimit(dto.Limit)
	}
	if dto.Offset > 0 {
		qb.SetOffset(dto.Offset)
	}

	query, args := qb.Build()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query search logs: %w", err)
	}
	defer rows.Close()

	logs := make([]*models.SearchLog, 0, 20)
	for rows.Next() {
		log := &models.SearchLog{}
		var queryBytes []byte

		if err := rows.Scan(
			&log.ID, &log.SearchId, &log.ActorID, &log.ActorName, &log.SearchType,
			&queryBytes, &log.DurationMs, &log.ResultsCount, &log.ItemsCount, &log.CreatedAt,
			&log.Actor.LastName, &log.Actor.FirstName, &log.Actor.Email,
		); err != nil {
			return nil, fmt.Errorf("failed to scan search log: %w", err)
		}

		if queryBytes != nil {
			log.Query = json.RawMessage(queryBytes)
		}

		logs = append(logs, log)
	}

	return logs, nil
}
