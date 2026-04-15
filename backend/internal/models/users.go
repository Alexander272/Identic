package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`

	AccessToken  string `json:"token"`
	RefreshToken string `json:"-"`
}

type Actor struct {
	ID   uuid.UUID
	Name string
}

type UserShort struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SSO_ID    string    `json:"ssoId" db:"sso_id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
}

type UserData struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SSO_ID    string    `json:"ssoId" db:"sso_id"`
	Role      string    `json:"role" db:"role"`
	RoleID    uuid.UUID `json:"roleId,omitempty" db:"role_id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	LastVisit time.Time `json:"lastVisit,omitempty" db:"last_visit"`
}

type UserDataDTO struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Actor     Actor
	SSO_ID    string    `json:"ssoId" db:"sso_id"`
	RoleID    uuid.UUID `json:"roleId" db:"role_id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	IsActive  bool      `json:"isActive" db:"is_active"`
}

type GetUserInfoDTO struct {
	UserID string `json:"userId"`
}

type UserRole struct {
	UserID   uuid.UUID
	RoleName string
}

type UserRoleDTO struct {
	UserID  uuid.UUID `json:"userId" db:"user_id"`
	RoleID  uuid.UUID `json:"roleId" db:"role_id"`
	ActorID Actor
}

// type GetByRealmDTO struct {
// 	RealmID string `json:"realmId" binding:"required"`
// 	Include bool   `json:"include"`
// }

// type GetByAccessDTO struct {
// 	RealmID string `json:"realmId" binding:"required"`
// 	Role    string `json:"role"`
// }

// type KeycloakUser struct {
// 	Id        string `json:"id"`
// 	Username  string `json:"username"`
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// 	Email     string `json:"email"`
// }
