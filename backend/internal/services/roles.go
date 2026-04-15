package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/events"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"golang.org/x/sync/errgroup"
)

type RoleService struct {
	repo      repository.Roles
	hierarchy RoleHierarchy
	perms     Permissions
	eventBus  *events.PolicyEventManager
}

type RoleDeps struct {
	Repo        repository.Roles
	Hierarchy   RoleHierarchy
	Permissions Permissions
	EventBus    *events.PolicyEventManager
}

func NewRoleService(deps *RoleDeps) *RoleService {
	return &RoleService{
		repo:      deps.Repo,
		hierarchy: deps.Hierarchy,
		perms:     deps.Permissions,
		eventBus:  deps.EventBus,
	}
}

type Roles interface {
	GetOne(ctx context.Context, req *models.GetRoleDTO) (*models.Role, error)
	GetAll(ctx context.Context) ([]*models.Role, error)
	GetWithStats(ctx context.Context) ([]*models.RoleWithStats, error)
	IsExists(ctx context.Context, roleName string) (bool, error)
	Create(ctx context.Context, dto *models.RoleDTO) error
	Update(ctx context.Context, dto *models.RoleDTO) error
	Delete(ctx context.Context, dto *models.DeleteRoleDTO) error
	AssignPermission(ctx context.Context, dto *models.RolePermissionDTO) error
	DeletePermission(ctx context.Context, dto *models.RolePermissionDTO) error
}

func (s *RoleService) GetOne(ctx context.Context, req *models.GetRoleDTO) (*models.Role, error) {
	data, err := s.repo.GetOne(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return data, nil
}

func (s *RoleService) GetAll(ctx context.Context) ([]*models.Role, error) {
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}
	return data, nil
}

func (s *RoleService) GetWithStats(ctx context.Context) ([]*models.RoleWithStats, error) {
	roles, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(roles))
	slugs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID.String())
		slugs = append(slugs, role.Slug)
	}

	var (
		userCounts  map[string]int
		permsCounts map[string]models.PermsCount
		inheritance map[string][]string
	)

	g, asyncCtx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		// Метод должен уметь возвращать дерево для списка ролей
		var err error
		inheritance, err = s.hierarchy.GetRoleDescendants(asyncCtx, &models.GetRolesInheritance{Roles: slugs})
		return err
	})
	g.Go(func() error {
		var err error
		userCounts, err = s.repo.GetUserCount(asyncCtx, roleIDs)
		if err != nil {
			return fmt.Errorf("failed to get user count: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if len(inheritance) > 0 {
		permsCounts, err = s.perms.CountForAll(ctx, inheritance)
		if err != nil {
			return nil, err
		}
	}

	result := make([]*models.RoleWithStats, 0, len(roles))
	for _, r := range roles {
		result = append(result, &models.RoleWithStats{
			Role:       *r,
			Inherited:  inheritance[r.Slug],
			UserCount:  userCounts[r.ID.String()],
			PermsCount: permsCounts[r.Slug],
		})
	}

	return result, nil
}

func (s *RoleService) IsExists(ctx context.Context, roleName string) (bool, error) {
	data, err := s.repo.IsExists(ctx, roleName)
	if err != nil {
		return false, fmt.Errorf("failed to check if role exists: %w", err)
	}
	return data, nil
}

func (s *RoleService) Create(ctx context.Context, dto *models.RoleDTO) error {
	err := s.repo.Create(ctx, nil, dto)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

func (s *RoleService) Update(ctx context.Context, dto *models.RoleDTO) error {
	err := s.repo.Update(ctx, nil, dto)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (s *RoleService) Delete(ctx context.Context, dto *models.DeleteRoleDTO) error {
	err := s.repo.Delete(ctx, nil, dto)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

func (s *RoleService) AssignPermission(ctx context.Context, dto *models.RolePermissionDTO) error {
	//TODO добавить транзакцию
	err := s.repo.AssignPermission(ctx, nil, dto)
	if err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}
	return nil
}

func (s *RoleService) DeletePermission(ctx context.Context, dto *models.RolePermissionDTO) error {
	//TODO добавить транзакцию
	err := s.repo.DeletePermission(ctx, nil, dto)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}
