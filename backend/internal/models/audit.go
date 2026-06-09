package models

import "github.com/google/uuid"

type UserAuditData struct {
	IsActive bool      `json:"isActive"`
	RoleID   uuid.UUID `json:"roleId"`
}
