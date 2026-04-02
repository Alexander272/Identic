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
	GetById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error)
	GetInfoById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error)
	GetByYear(ctx context.Context, req *models.GetOrderByYearDTO) ([]*models.Order, error)
	GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error)
	GetFlatData(ctx context.Context, req *models.GetFlatOrderDTO) (*models.FlatOrderRes, error)
	Create(ctx context.Context, dto *models.OrderDTO) error
	CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.OrderDTO) error
	Update(ctx context.Context, dto *models.OrderDTO) error
}

func (s *OrdersService) GetById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error) {
	data, err := s.repo.GetById(ctx, req)
	if err != nil {
		if err == models.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get order. error: %w", err)
	}
	positions, err := s.positions.GetByOrder(ctx, &models.GetPositionsByOrderIdDTO{OrderId: data.Id})
	if err != nil {
		return nil, err
	}
	data.Positions = positions

	return data, nil
}

func (s *OrdersService) GetInfoById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error) {
	data, err := s.repo.GetById(ctx, req)
	if err != nil {
		if err == models.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get order info. error: %w", err)
	}

	if len(req.PositionIds) > 0 {
		positions, err := s.positions.GetByIds(ctx, &models.GetPositionsByIds{Ids: req.PositionIds})
		if err != nil {
			return nil, err
		}
		data.Positions = positions
	}

	return data, nil
}

func (s *OrdersService) GetByYear(ctx context.Context, req *models.GetOrderByYearDTO) ([]*models.Order, error) {
	data, err := s.repo.GetByYear(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders. error: %w", err)
	}
	return data, nil
}

func (s *OrdersService) GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error) {
	data, err := s.repo.GetUniqueData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique data. error: %w", err)
	}
	return data, nil
}

func (s *OrdersService) GetFlatData(ctx context.Context, req *models.GetFlatOrderDTO) (*models.FlatOrderRes, error) {
	data, err := s.repo.GetFlatData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get flat data. error: %w", err)
	}
	return data, nil
}

func (s *OrdersService) IsExist(ctx context.Context, tx postgres.Tx, dto *models.OrderDTO) (bool, error) {
	exist, err := s.repo.IsExist(ctx, tx, dto)
	if err != nil {
		return false, fmt.Errorf("failed to check if order exists. error: %w", err)
	}
	return exist, nil
}

func (s *OrdersService) Create(ctx context.Context, dto *models.OrderDTO) error {
	return s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		isExist, err := s.repo.IsExist(ctx, tx, dto)
		if err != nil {
			return err
		}

		if isExist {
			return models.ErrOrderAlreadyExists
		}

		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create order. error: %w", err)
		}
		for i := range dto.Positions {
			dto.Positions[i].OrderId = dto.Id
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

		for i := range dto.Positions {
			dto.Positions[i].Id = uuid.NewString()
			dto.Positions[i].OrderId = dto.Id
		}

		if err := s.positions.DeleteByOrder(ctx, tx, &models.DeletePositionsByOrderIdDTO{OrderId: dto.Id}); err != nil {
			return err
		}
		if err := s.positions.Create(ctx, tx, dto.Positions); err != nil {
			return err
		}
		return nil
	})
}
