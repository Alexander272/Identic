package models

import "github.com/google/uuid"

type User struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`

	AccessToken  string `json:"token"`
	RefreshToken string `json:"-"`
}

type UserData struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SSO_ID    string    `json:"ssoId" db:"sso_id"`
	RoleID    uuid.UUID `json:"roleId" db:"role_id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Email     string    `json:"email" db:"email"`
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
	ActorID uuid.UUID `json:"actorId" db:"actor_id"`
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
