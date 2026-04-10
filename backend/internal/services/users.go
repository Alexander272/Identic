package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/pkg/auth"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/Nerzal/gocloak/v13"
)

type userService struct {
	repo     repository.Users
	tm       TransactionManager
	keycloak *auth.KeycloakClient
	role     Roles
}

type UsersDeps struct {
	Repo      repository.Users
	TxManager TransactionManager
	Keycloak  *auth.KeycloakClient
	Role      Roles
}

func NewUserService(deps *UsersDeps) *userService {
	return &userService{
		repo:     deps.Repo,
		tm:       deps.TxManager,
		keycloak: deps.Keycloak,
		role:     deps.Role,
	}
}

type Users interface {
	LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.UserRole, error)
	GetAll(ctx context.Context) ([]*models.UserData, error)
	Sync(ctx context.Context) error
}

func (s *userService) LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.UserRole, error) {
	data, err := s.repo.LoadPolicy(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}
	return data, nil
}

func (s *userService) GetAll(ctx context.Context) ([]*models.UserData, error) {
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users. error: %w", err)
	}
	return data, nil
}

func (s *userService) Sync(ctx context.Context) error {
	logger.Info("Sync users started")

	token, err := s.keycloak.Login(ctx)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	// 1. Быстрый поиск ID группы
	groups, err := s.keycloak.Client.GetGroups(ctx, token.AccessToken, s.keycloak.Realm, gocloak.GetGroupsParams{
		Search: gocloak.StringP("identic"), // Фильтруем на стороне Keycloak
	})
	if err != nil {
		return fmt.Errorf("failed to get groups: %w", err)
	}

	var groupID string
	for _, g := range groups {
		if g.Name != nil && *g.Name == "identic" {
			groupID = *g.ID
			break
		}
	}
	if groupID == "" {
		return fmt.Errorf("group 'identic' not found")
	}

	// 2. Получаем активных пользователей из Keycloak
	keycloakUsers, err := s.keycloak.Client.GetGroupMembers(ctx, token.AccessToken, s.keycloak.Realm, groupID, gocloak.GetGroupsParams{Max: gocloak.IntP(1000)})
	if err != nil {
		return fmt.Errorf("failed to get group members: %w", err)
	}

	if len(keycloakUsers) == 0 {
		return fmt.Errorf("group 'identic' is empty")
	}

	// Пред-аллокация для пачки данных из Keycloak
	kcDataMap := make(map[string]*models.UserData, len(keycloakUsers))
	for _, u := range keycloakUsers {
		if u.Enabled != nil && !*u.Enabled {
			continue
		}

		userData := s.mapToUserData(u)
		kcDataMap[userData.SSO_ID] = userData
	}

	// 3. Получаем текущих пользователей из нашей БД
	dbUsers, err := s.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch DB users: %w", err)
	}

	defRole, err := s.role.GetOne(ctx, &models.GetRoleDTO{Slug: "user"})
	if err != nil {
		return err
	}

	toCreate := make([]*models.UserData, 0)
	toUpdate := make([]*models.UserData, 0)
	toDelete := make([]string, 0)

	// 4. Основной цикл синхронизации
	dbUserMap := make(map[string]*models.UserData, len(dbUsers))

	for _, dbU := range dbUsers {
		dbUserMap[dbU.SSO_ID] = dbU

		if kcData, exists := kcDataMap[dbU.SSO_ID]; exists {
			// Проверяем, нужно ли реально обновлять (DeepEqual или по полям)
			if s.isChanged(dbU, kcData) {
				toUpdate = append(toUpdate, kcData)
			}
			// Удаляем из мапы Keycloak, чтобы там остались только "новые"
			delete(kcDataMap, dbU.SSO_ID)
		} else {
			// Если в Keycloak нет, а в БД есть — на удаление
			toDelete = append(toDelete, dbU.SSO_ID)
		}
	}

	// Все, кто остались в kcDataMap — новые
	for _, newU := range kcDataMap {
		newU.RoleID = defRole.ID
		toCreate = append(toCreate, newU)
	}

	// 5. Выполнение операций (Batch processing)
	return s.tm.WithinTransaction(ctx, func(tx postgres.Tx) error {
		if len(toCreate) > 0 {
			if err := s.CreateSeveral(ctx, tx, toCreate); err != nil {
				return err
			}
		}
		if len(toUpdate) > 0 {
			if err := s.UpdateSeveral(ctx, tx, toUpdate); err != nil {
				return err
			}
		}
		if len(toDelete) > 0 {
			if err := s.DeleteSeveral(ctx, tx, toDelete); err != nil {
				return err
			}
		}

		logger.Info("Sync finished",
			"created", len(toCreate),
			"updated", len(toUpdate),
			"deleted", len(toDelete))
		return nil
	})
}

// Вспомогательная функция для маппинга (убирает дублирование nil-проверок)
func (s *userService) mapToUserData(u *gocloak.User) *models.UserData {
	return &models.UserData{
		SSO_ID:    s.nonNil(u.ID),
		Username:  s.nonNil(u.Username),
		Email:     s.nonNil(u.Email),
		FirstName: s.nonNil(u.FirstName),
		LastName:  s.nonNil(u.LastName),
	}
}

func (s *userService) nonNil(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// Функция проверки изменений, чтобы не дергать БД зря
func (s *userService) isChanged(old, new *models.UserData) bool {
	return old.Username != new.Username ||
		old.Email != new.Email ||
		old.FirstName != new.FirstName ||
		old.LastName != new.LastName
}

func (s *userService) CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.UserData) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.CreateSeveral(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create few users. error: %w", err)
	}
	return nil
}

func (s *userService) UpdateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.UserData) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.UpdateSeveral(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to update few users. error: %w", err)
	}
	return nil
}

func (s *userService) DeleteSeveral(ctx context.Context, tx postgres.Tx, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.DeleteSeveral(ctx, tx, ids); err != nil {
		return fmt.Errorf("failed to delete few users. error: %w", err)
	}
	return nil
}
