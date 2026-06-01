package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/goccy/go-json"
)

type PriceSearchLogRepo struct {
	db *pgxpool.Pool
}

func NewPriceSearchLogRepo(db *pgxpool.Pool) *PriceSearchLogRepo {
	return &PriceSearchLogRepo{db: db}
}

type PriceSearchLogs interface {
	Create(ctx context.Context, dto *models.CreatePriceSearchLogDTO) error
	Get(ctx context.Context, dto *models.GetPriceSearchLogsDTO) ([]*models.PriceSearchLog, error)
}

func (r *PriceSearchLogRepo) Create(ctx context.Context, dto *models.CreatePriceSearchLogDTO) error {
	queriesJSON, err := json.Marshal(dto.Queries)
	if err != nil {
		return fmt.Errorf("failed to marshal queries: %w", err)
	}

	codesJSON, err := json.Marshal(dto.Codes)
	if err != nil {
		return fmt.Errorf("failed to marshal codes: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (queries, codes, actor_id, actor_name, results_count, duration_ms)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		Tables.PriceSearchLogs,
	)

	_, err = r.db.Exec(ctx, query,
		queriesJSON,
		codesJSON,
		dto.ActorID,
		dto.ActorName,
		dto.ResultsCount,
		dto.DurationMs,
	)
	if err != nil {
		return postgres.MapError(fmt.Errorf("failed to create price search log: %w", err))
	}

	return nil
}

func (r *PriceSearchLogRepo) Get(ctx context.Context, dto *models.GetPriceSearchLogsDTO) ([]*models.PriceSearchLog, error) {
	baseQuery := fmt.Sprintf(`SELECT id, queries, codes, actor_id, actor_name, results_count, duration_ms, created_at FROM %s`,
		Tables.PriceSearchLogs,
	)

	qb := postgres.NewQueryBuilder(baseQuery)
	qb.AddUUIDFilter("actor_id", dto.ActorID)
	qb.AddDateRangeFilter("created_at", dto.StartDate, dto.EndDate)
	qb.SetSort("created_at", true)

	if dto.Limit > 0 {
		qb.SetLimit(dto.Limit)
	}
	if dto.Offset > 0 {
		qb.SetOffset(dto.Offset)
	}

	query, args := qb.Build()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, postgres.MapError(fmt.Errorf("failed to query price search logs: %w", err))
	}
	defer rows.Close()

	logs := make([]*models.PriceSearchLog, 0, 20)
	for rows.Next() {
		log := &models.PriceSearchLog{}
		var queriesBytes, codesBytes []byte

		if err := rows.Scan(
			&log.ID, &queriesBytes, &codesBytes, &log.ActorID, &log.ActorName,
			&log.ResultsCount, &log.DurationMs, &log.CreatedAt,
		); err != nil {
			return nil, postgres.MapError(fmt.Errorf("failed to scan price search log: %w", err))
		}

		if queriesBytes != nil {
			log.Queries = json.RawMessage(queriesBytes)
		}
		if codesBytes != nil {
			json.Unmarshal(codesBytes, &log.Codes)
		}

		logs = append(logs, log)
	}

	return logs, nil
}
