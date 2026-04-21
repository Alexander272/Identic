package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
)

type userLoginService struct {
	repo repository.UserLogins
	tm   TransactionManager
}

func NewUserLoginService(repo repository.UserLogins, tm TransactionManager) *userLoginService {
	return &userLoginService{
		repo: repo,
		tm:   tm,
	}
}

type UserLogins interface {
	RecordLogin(ctx context.Context, dto *models.UserLoginDTO) error
	GetByUser(ctx context.Context, req *models.GetUserLoginsDTO) ([]*models.UserLogin, int64, error)
	GetLastByUser(ctx context.Context, userID string) (*models.UserLogin, error)
	GetLastByUsers(ctx context.Context, req *models.GetUserLoginsDTO) ([]*models.UserLoginWithUser, error)
	UpdateLastActivity(ctx context.Context, tx postgres.Tx, userID string) (bool, error)
}

func (s *userLoginService) RecordLogin(ctx context.Context, dto *models.UserLoginDTO) error {
	if dto.Metadata == nil {
		metadata := models.LoginMetadata{Success: true}
		data, _ := json.Marshal(metadata)
		dto.Metadata = data
	}
	if err := s.repo.Create(ctx, nil, dto); err != nil {
		return fmt.Errorf("failed to create user login: %w", err)
	}
	return nil
}

func (s *userLoginService) GetByUser(ctx context.Context, req *models.GetUserLoginsDTO) ([]*models.UserLogin, int64, error) {
	if req.Limit == 0 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	data, err := s.repo.GetByUser(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user logins: %w", err)
	}

	count, err := s.repo.GetByUserCount(ctx, req.UserID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user logins: %w", err)
	}

	return data, count, nil
}

func (s *userLoginService) GetLastByUser(ctx context.Context, userID string) (*models.UserLogin, error) {
	login, err := s.repo.GetLastByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last user login: %w", err)
	}
	return login, nil
}

func (s *userLoginService) GetLastByUsers(ctx context.Context, req *models.GetUserLoginsDTO) ([]*models.UserLoginWithUser, error) {
	logins, err := s.repo.GetLastByUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get last user logins: %w", err)
	}
	return logins, nil
}

func (s *userLoginService) UpdateLastActivity(ctx context.Context, tx postgres.Tx, userID string) (bool, error) {
	wasIdle, err := s.repo.UpdateLastActivity(ctx, tx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to update last activity: %w", err)
	}
	return wasIdle, nil
}
