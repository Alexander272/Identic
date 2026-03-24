package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
)

type PositionsService struct {
	repo      repository.Positions
	txManager TransactionManager
}

func NewPositionsService(repo repository.Positions, txManager TransactionManager) *PositionsService {
	return &PositionsService{
		repo:      repo,
		txManager: txManager,
	}
}

type Positions interface {
	GetByOrder(ctx context.Context, req *models.GetPositionsByOrderIdDTO) ([]*models.Position, error)
	GetByIds(ctx context.Context, req *models.GetPositionsByIds) ([]*models.Position, error)
	Create(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error
	Update(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error
	Delete(ctx context.Context, tx postgres.Tx, dto []*models.DeletePositionDTO) error
}

func (s *PositionsService) GetByOrder(ctx context.Context, req *models.GetPositionsByOrderIdDTO) ([]*models.Position, error) {
	data, err := s.repo.GetByOrder(ctx, req)
	if err != nil {
		if err == models.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get positions. error: %w", err)
	}
	return data, nil
}

func (s *PositionsService) GetByIds(ctx context.Context, req *models.GetPositionsByIds) ([]*models.Position, error) {
	data, err := s.repo.GetByIds(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions. error: %w", err)
	}
	return data, nil
}

func (s *PositionsService) Create(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	if tx == nil {
		// Если транзакция не передана, создаем новую
		return s.txManager.WithinTransaction(ctx, func(newTx postgres.Tx) error {
			return s.executeCreate(ctx, newTx, dto)
		})
	}
	// Если транзакция передана, используем её
	return s.executeCreate(ctx, tx, dto)
}
func (s *PositionsService) executeCreate(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error {
	for i := range dto {
		dto[i].Search = NormalizeString(dto[i].Name)
	}

	if err := s.repo.Create(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create positions. error: %w", err)
	}
	return nil
}

func (s *PositionsService) Update(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	if tx == nil {
		// Если транзакция не передана, создаем новую
		return s.txManager.WithinTransaction(ctx, func(newTx postgres.Tx) error {
			return s.executeUpdate(ctx, newTx, dto)
		})
	}
	// Если транзакция передана, используем её
	return s.executeUpdate(ctx, tx, dto)
}
func (s *PositionsService) executeUpdate(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error {
	for i := range dto {
		dto[i].Search = NormalizeString(dto[i].Search)
	}

	if err := s.repo.Update(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to update positions. error: %w", err)
	}
	return nil
}

func (s *PositionsService) Delete(ctx context.Context, tx postgres.Tx, dto []*models.DeletePositionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	if tx == nil { // Если транзакция не передана, создаем новую
		return s.txManager.WithinTransaction(ctx, func(newTx postgres.Tx) error {
			return s.executeDelete(ctx, newTx, dto)
		})
	}
	return s.executeDelete(ctx, tx, dto) // Если транзакция передана, используем её
}
func (s *PositionsService) executeDelete(ctx context.Context, tx postgres.Tx, dto []*models.DeletePositionDTO) error {
	if err := s.repo.Delete(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to delete positions. error: %w", err)
	}
	return nil
}
