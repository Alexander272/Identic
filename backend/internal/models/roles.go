package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Level       int       `json:"level" db:"level"`
	IsActive    bool      `json:"isActive" db:"is_active"`
	IsSystem    bool      `json:"isSystem" db:"is_system"`
	IsEditable  bool      `json:"isEditable" db:"is_editable"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type RoleWithStats struct {
	Role
	Inherited  []string   `json:"inherited"`
	PermsCount PermsCount `json:"perms"`
	UserCount  int        `json:"userCount"`
}

type RoleShort struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Slug string    `json:"slug" db:"slug"`
	Name string    `json:"name" db:"name"`
}

type GetRoleDTO struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Slug string    `json:"slug" db:"slug"`
}

type RoleDTO struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Actor       Actor
	Slug        string    `json:"slug" db:"slug"`
	Name        string    `json:"name" db:"name"`
	Level       int       `json:"level" db:"level"`
	IsSystem    bool      `json:"isSystem" db:"is_system"`
	Permissions []string  `json:"permissions" db:"permissions"`
	Children    []string  `json:"children" db:"children"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type DeleteRoleDTO struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Actor Actor
}

type RoleInheritance struct {
	ParentRole string
	ChildRole  string
	Realm      string
}

type RolePermission struct {
	RoleID       uuid.UUID `json:"roleId" db:"role_id"`
	PermissionID uuid.UUID `json:"permissionId" db:"permission_id"`
}

type RolePermissionDTO struct {
	ActorID      uuid.UUID `json:"actorId" db:"actor_id"`
	RoleID       uuid.UUID `json:"roleId" db:"role_id"`
	PermissionID uuid.UUID `json:"permissionId" db:"permission_id"`
}
