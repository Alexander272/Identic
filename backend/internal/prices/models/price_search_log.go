package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PriceSearchLog struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	Queries      json.RawMessage `json:"queries" db:"queries"`
	Codes        []string        `json:"codes,omitempty" db:"codes"`
	ActorID      uuid.UUID       `json:"actorId" db:"actor_id"`
	ActorName    string          `json:"actorName" db:"actor_name"`
	ResultsCount int             `json:"resultsCount" db:"results_count"`
	DurationMs   int64           `json:"durationMs" db:"duration_ms"`
	CreatedAt    time.Time       `json:"createdAt" db:"created_at"`
}

type CreatePriceSearchLogDTO struct {
	Queries      []string  `json:"queries"`
	Codes        []string  `json:"codes"`
	ActorID      uuid.UUID `json:"actorId"`
	ActorName    string    `json:"actorName"`
	ResultsCount int       `json:"resultsCount"`
	DurationMs   int64     `json:"durationMs"`
}

type GetPriceSearchLogsDTO struct {
	ActorID   *uuid.UUID `json:"actorId,omitempty"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
}
