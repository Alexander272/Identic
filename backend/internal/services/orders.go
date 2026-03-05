package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/google/uuid"
)

type OrdersService struct {
	repo      repository.Orders
	txManager TransactionManager
	positions Positions
}

func NewOrdersService(repo repository.Orders, txManager TransactionManager, positions Positions) *OrdersService {
	return &OrdersService{
		repo:      repo,
		txManager: txManager,
		positions: positions,
	}
}

type Orders interface {
	Create(ctx context.Context, dto *models.OrderDTO) error
	CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.OrderDTO) error
	Update(ctx context.Context, dto *models.OrderDTO) error
}

func (s *OrdersService) Create(ctx context.Context, dto *models.OrderDTO) error {
	return s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create order. error: %w", err)
		}
		if err := s.positions.Create(ctx, tx, dto.Positions); err != nil {
			return fmt.Errorf("failed to create positions. error: %w", err)
		}
		return nil
	})
}

func (s *OrdersService) CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.OrderDTO) error {
	if tx == nil {
		// Если транзакция не передана, создаем новую
		return s.txManager.WithinTransaction(ctx, func(newTx postgres.Tx) error {
			return s.executeCreate(ctx, newTx, dto)
		})
	}
	// Если транзакция передана, используем её
	return s.executeCreate(ctx, tx, dto)
}
func (s *OrdersService) executeCreate(ctx context.Context, tx postgres.Tx, dto []*models.OrderDTO) error {
	positions := make([]*models.PositionDTO, 0, len(dto))

	for i := range dto {
		dto[i].Id = uuid.NewString()

		for j := range dto[i].Positions {
			dto[i].Positions[j].Id = uuid.NewString()
			dto[i].Positions[j].OrderId = dto[i].Id

			positions = append(positions, dto[i].Positions[j])
		}
	}

	if err := s.repo.CreateSeveral(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create orders. error: %w", err)
	}
	if err := s.positions.Create(ctx, tx, positions); err != nil {
		return fmt.Errorf("failed to create positions. error: %w", err)
	}
	return nil
}

func (s *OrdersService) Update(ctx context.Context, dto *models.OrderDTO) error {
	return s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update order. error: %w", err)
		}
		if err := s.positions.Update(ctx, tx, dto.Positions); err != nil {
			return fmt.Errorf("failed to update positions. error: %w", err)
		}
		return nil
	})
}
