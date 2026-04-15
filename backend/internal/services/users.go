package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/events"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/pkg/auth"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/Nerzal/gocloak/v13"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type userService struct {
	repo     repository.Users
	tm       TransactionManager
	keycloak *auth.KeycloakClient
	role     Roles
	eventBus *events.PolicyEventManager
}

type UsersDeps struct {
	Repo      repository.Users
	TxManager TransactionManager
	Keycloak  *auth.KeycloakClient
	Role      Roles
	EventBus  *events.PolicyEventManager
}

func NewUserService(deps *UsersDeps) *userService {
	return &userService{
		repo:     deps.Repo,
		tm:       deps.TxManager,
		keycloak: deps.Keycloak,
		role:     deps.Role,
		eventBus: deps.EventBus,
	}
}

type Users interface {
	LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.UserRole, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.UserData, error)
	GetAll(ctx context.Context) ([]*models.UserData, error)
	Sync(ctx context.Context, actor *models.Actor) error
	Update(ctx context.Context, dto *models.UserDataDTO) error
}

func (s *userService) LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.UserRole, error) {
	data, err := s.repo.LoadPolicy(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}
	return data, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.UserData, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id. error: %w", err)
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

func (s *userService) Sync(ctx context.Context, actor *models.Actor) error {
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
	kcDataMap := make(map[string]*models.UserDataDTO, len(keycloakUsers))
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

	toCreate := make([]*models.UserDataDTO, 0)
	toUpdate := make([]*models.UserDataDTO, 0)
	toDelete := make([]string, 0)

	// 4. Основной цикл синхронизации
	for _, dbU := range dbUsers {
		existUser := &models.UserDataDTO{
			ID:        dbU.ID,
			SSO_ID:    dbU.SSO_ID,
			Username:  dbU.Username,
			FirstName: dbU.FirstName,
			LastName:  dbU.LastName,
			Email:     dbU.Email,
		}

		if kcData, exists := kcDataMap[dbU.SSO_ID]; exists {
			// Проверяем, нужно ли реально обновлять (DeepEqual или по полям)
			if s.isChanged(existUser, kcData) {
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
	err = s.tm.WithinTransaction(ctx, func(tx postgres.Tx) error {
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

	if err != nil {
		return fmt.Errorf("failed to execute batch: %w", err)
	}

	// Формируем событие для Casbin и Audit
	event := events.PolicyEvent{
		ChangedBy:     actor.ID,
		ChangedByName: actor.Name,
		Action:        "sync_users",
		EntityType:    "users",
	}
	// Отправляем в шину
	s.eventBus.Notify(event)
	return nil
}

// Вспомогательная функция для маппинга (убирает дублирование nil-проверок)
func (s *userService) mapToUserData(u *gocloak.User) *models.UserDataDTO {
	return &models.UserDataDTO{
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
func (s *userService) isChanged(old, new *models.UserDataDTO) bool {
	return old.Username != new.Username ||
		old.Email != new.Email ||
		old.FirstName != new.FirstName ||
		old.LastName != new.LastName
}

func (s *userService) CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.UserDataDTO) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.CreateSeveral(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create few users. error: %w", err)
	}
	return nil
}

func (s *userService) Update(ctx context.Context, dto *models.UserDataDTO) error {
	candidate, err := s.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, nil, dto); err != nil {
		return fmt.Errorf("failed to update user. error: %w", err)
	}

	oldValue, err := json.Marshal(models.UserAuditData{
		IsActive: candidate.IsActive,
		RoleID:   candidate.RoleID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal policy event. error: %w", err)
	}

	newValue, err := json.Marshal(models.UserAuditData{
		IsActive: dto.IsActive,
		RoleID:   dto.RoleID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal policy event. error: %w", err)
	}

	// Формируем событие для Casbin и Audit
	event := events.PolicyEvent{
		ChangedBy:     dto.Actor.ID,
		ChangedByName: dto.Actor.Name,
		Action:        "update_user",
		EntityType:    "users",
		EntityID:      &dto.ID,
		OldValues:     oldValue,
		NewValues:     newValue,
	}
	// Отправляем в шину
	s.eventBus.Notify(event)

	return nil
}

func (s *userService) UpdateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.UserDataDTO) error {
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
