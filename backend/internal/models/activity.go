package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ActionType string

const (
	ActionInsert ActionType = "INSERT"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type EntityType string

const (
	EntityOrder     EntityType = "order"
	EntityOrderItem EntityType = "order_item"
)

type ActivityLog struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	Action        ActionType      `json:"action" db:"action"`
	ChangedBy     uuid.UUID       `json:"changedBy" db:"changed_by"`
	ChangedByName string          `json:"changedByName" db:"changed_by_name"`
	EntityType    EntityType      `json:"entityType" db:"entity_type"`
	EntityID      string          `json:"entityId" db:"entity_id"`
	Entity        *string         `json:"entity,omitempty" db:"entity"`
	ParentID      *string         `json:"parentId,omitempty" db:"parent_id"`
	OldValues     json.RawMessage `json:"oldValues,omitempty" db:"old_values"`
	NewValues     json.RawMessage `json:"newValues,omitempty" db:"new_values"`
	CreatedAt     time.Time       `json:"createdAt" db:"created_at"`
}

type CreateActivityLogDTO struct {
	Action        ActionType  `json:"action"`
	ChangedBy     uuid.UUID   `json:"changedBy"`
	ChangedByName string      `json:"changedByName"`
	EntityType    EntityType  `json:"entityType"`
	EntityID      string      `json:"entityId"`
	Entity        string      `json:"entity"` // Название объекта (напр. "Заказ #123: Иванов")
	ParentID      *string     `json:"parentId,omitempty"`
	OldValues     interface{} `json:"oldValues,omitempty"`
	NewValues     interface{} `json:"newValues,omitempty"`
}

type GetActivityLogsDTO struct {
	EntityID   string     `json:"entityId"`
	EntityType EntityType `json:"entityType,omitempty"`
	ParentID   *uuid.UUID `json:"parentId,omitempty"`
}

type GetAllActivityLogsDTO struct {
	ActorID   *uuid.UUID `json:"actorId,omitempty"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
}

type BatchLogPositionsDTO struct {
	OrderID string
	Actor   Actor
	Created []*PositionDTO
	Updated []*PositionDTO
	Deleted []*PositionDTO
	Old     []*Position
}

type OrderDiff struct {
	OldValues map[string]interface{} `json:"oldValues,omitempty"`
	NewValues map[string]interface{} `json:"newValues,omitempty"`
}

type OrderLogMode string

const (
	OrderLogDiff OrderLogMode = "diff" // только изменённые поля
	OrderLogFull OrderLogMode = "full" // полные снапшоты
)
