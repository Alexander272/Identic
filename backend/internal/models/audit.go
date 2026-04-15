package models

import "github.com/google/uuid"

type UserAuditData struct {
	IsActive bool      `json:"is_active"`
	RoleID   uuid.UUID `json:"role_id"`
}
