package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/events"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/google/uuid"
)

type PermissionService struct {
	repo     repository.Permissions
	tm       TransactionManager
	eventBus *events.PolicyEventManager
}

func NewPermissionService(repo repository.Permissions, tm TransactionManager, eventBus *events.PolicyEventManager) *PermissionService {
	return &PermissionService{
		repo:     repo,
		tm:       tm,
		eventBus: eventBus,
	}
}

type Permissions interface {
	GetByRole(ctx context.Context, req *models.GetPermsByRoleDTO) ([]*models.Permission, error)
	LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.Permission, error)
	Create(ctx context.Context, tx postgres.Tx, dto *models.PermissionDTO) error
	Delete(ctx context.Context, tx postgres.Tx, dto *models.DeletePermissionDTO) error
}

func (s *PermissionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	data, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission by id: %w", err)
	}
	return data, nil
}

func (s *PermissionService) GetByRole(ctx context.Context, req *models.GetPermsByRoleDTO) ([]*models.Permission, error) {
	data, err := s.repo.GetByRole(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by role: %w", err)
	}
	return data, nil
}

func (s *PermissionService) LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.Permission, error) {
	data, err := s.repo.LoadPolicy(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}
	return data, nil
}

func (s *PermissionService) Create(ctx context.Context, tx postgres.Tx, dto *models.PermissionDTO) error {
	// if constants.ResourcesList.Permissions

	err := s.repo.Create(ctx, tx, dto)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	return nil
}

func (s *PermissionService) Delete(ctx context.Context, tx postgres.Tx, dto *models.DeletePermissionDTO) error {
	err := s.repo.Delete(ctx, tx, dto)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}
