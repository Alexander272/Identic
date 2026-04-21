package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SearchType string

const (
	SearchTypeExact SearchType = "exact"
	SearchTypeFuzzy SearchType = "fuzzy"
)

type SearchLog struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	SearchId     string          `json:"searchId" db:"search_id"`
	Actor        UserShort       `json:"actor" db:"actor"`
	ActorID      uuid.UUID       `json:"actorId" db:"actor_id"`
	ActorName    string          `json:"actorName" db:"actor_name"`
	SearchType   SearchType      `json:"searchType" db:"search_type"`
	Query        json.RawMessage `json:"query" db:"query"`
	DurationMs   int64           `json:"durationMs" db:"duration_ms"`
	ResultsCount int             `json:"resultsCount" db:"results_count"`
	ItemsCount   int             `json:"itemsCount" db:"items_count"`
	CreatedAt    time.Time       `json:"createdAt" db:"created_at"`
}

type CreateSearchLogDTO struct {
	SearchId     string      `json:"searchId"`
	ActorID      uuid.UUID   `json:"actorId"`
	ActorName    string      `json:"actorName"`
	SearchType   SearchType  `json:"searchType"`
	Query        interface{} `json:"query"`
	DurationMs   int64       `json:"durationMs"`
	ResultsCount int         `json:"resultsCount"`
	ItemsCount   int         `json:"itemsCount"`
}

type GetSearchLogsDTO struct {
	ActorID   *uuid.UUID `json:"actorId,omitempty"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
}
