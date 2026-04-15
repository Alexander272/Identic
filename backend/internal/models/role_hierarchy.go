package models

import "github.com/google/uuid"

type RoleHierarchy struct {
	Role       Role `json:"childRole"`
	ParentRole Role `json:"parentRole"`
}

type RoleWithHierarchy struct {
	Role
	InheritsFrom []*RoleHierarchy `json:"inherits_from"` // от кого наследуем
	InheritedBy  []*RoleHierarchy `json:"inherited_by"`  // кто наследует от нас
}

type RoleHierarchyDTO struct {
	ParentRoleID uuid.UUID `json:"parentRoleId" db:"parent_role_id"`
	RoleID       uuid.UUID `json:"childRoleId" db:"child_role_id"`
	ActorID      uuid.UUID `json:"actorId" db:"actor_id"`
}

type GetRoleInheritance struct {
	Role string
}

type GetRolesInheritance struct {
	Roles []string
}

type SyncRoleInheritance struct {
	Role       string
	ParentRole string
}
