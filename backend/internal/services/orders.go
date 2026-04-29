package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/google/uuid"
)

type OrdersService struct {
	repo      repository.Orders
	txManager TransactionManager
	positions Positions
	search    Search
	activity  Activity
}

func NewOrdersService(repo repository.Orders, txManager TransactionManager, positions Positions, search Search, activity Activity) *OrdersService {
	return &OrdersService{
		repo:      repo,
		txManager: txManager,
		positions: positions,
		search:    search,
		activity:  activity,
	}
}

type Orders interface {
	Get(ctx context.Context, req *models.OrderFilterDTO) ([]*models.Order, error)
	GetById(ctx context.Context, tx postgres.Tx, req *models.GetOrderByIdDTO) (*models.Order, error)
	GetInfoById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error)
	GetByYear(ctx context.Context, req *models.GetOrderByYearDTO) ([]*models.Order, error)
	GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error)
	GetFlatData(ctx context.Context, req *models.GetFlatOrderDTO) (*models.FlatOrderRes, error)
	Create(ctx context.Context, dto *models.OrderDTO) (string, error)
	CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.OrderDTO) error
	Update(ctx context.Context, dto *models.OrderDTO) error
}

func (s *OrdersService) Get(ctx context.Context, req *models.OrderFilterDTO) ([]*models.Order, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders. error: %w", err)
	}
	return data, nil
}

func (s *OrdersService) GetById(ctx context.Context, tx postgres.Tx, req *models.GetOrderByIdDTO) (*models.Order, error) {
	data, err := s.repo.GetById(ctx, tx, req)
	if err != nil {
		if err == models.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get order. error: %w", err)
	}
	positions, err := s.positions.GetByOrder(ctx, tx, &models.GetPositionsByOrderIdDTO{OrderId: data.Id})
	if err != nil {
		return nil, err
	}
	data.Positions = positions

	if req.SearchId != "" {
		cache := &models.GetCacheDTO{OrderId: req.Id, SearchId: req.SearchId}
		posIds, err := s.search.GetCache(ctx, cache)
		if err != nil {
			return nil, err
		}

		found := make(map[string]struct{}, len(posIds))
		for _, posId := range posIds {
			found[posId] = struct{}{}
		}
		for _, pos := range data.Positions {
			if _, ok := found[pos.Id]; ok {
				pos.IsFound = true
			}
		}
		data.PosWereFound = len(found) > 0
	}

	return data, nil
}

func (s *OrdersService) GetInfoById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error) {
	data, err := s.repo.GetById(ctx, nil, req)
	if err != nil {
		if err == models.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get order info. error: %w", err)
	}

	if req.SearchId != "" {
		cache := &models.GetCacheDTO{OrderId: req.Id, SearchId: req.SearchId}
		posIds, err := s.search.GetCache(ctx, cache)
		if err != nil {
			return nil, err
		}

		if len(posIds) > 0 {
			positions, err := s.positions.GetByIds(ctx, &models.GetPositionsByIds{Ids: posIds})
			if err != nil {
				return nil, err
			}
			data.Positions = positions
		}
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

func (s *OrdersService) Create(ctx context.Context, dto *models.OrderDTO) (string, error) {
	for i := range dto.Positions {
		dto.Positions[i].Name = ClearString(dto.Positions[i].Name)
	}

	dto.Hash = CalculateHash(dto.Positions)
	logger.Debug("create", logger.StringAttr("hash", dto.Hash), logger.IntAttr("len", len(dto.Positions)))

	err := s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		existingId, err := s.repo.IsExistByPos(ctx, tx, dto)
		if err != nil {
			return fmt.Errorf("failed to check if order exists. error: %w", err)
		}

		if existingId != "" {
			// Заказ уже существует - возвращаем его ID
			dto.Id = existingId
			return models.ErrOrderAlreadyExists
		}

		if dto.Id == "" {
			dto.Id = uuid.NewString()
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
	if err != nil {
		return dto.Id, err
	}

	go s.activity.AsyncLog(context.Background(), func() error {
		return s.txManager.WithinTransaction(context.Background(), func(tx postgres.Tx) error {
			if err := s.activity.LogOrderCreate(context.Background(), tx, dto); err != nil {
				return fmt.Errorf("failed to log order create: %w", err)
			}
			return s.activity.BatchLogPositions(context.Background(), tx, &models.BatchLogPositionsDTO{
				OrderID: dto.Id,
				Actor:   dto.Actor,
				Created: dto.Positions,
			})
		})
	}, map[string]any{
		"order_id": dto.Id,
		"action":   "order_" + models.ActionInsert,
		"actor":    dto.Actor,
	})

	return dto.Id, nil
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
	for i := range dto.Positions {
		dto.Positions[i].Name = ClearString(dto.Positions[i].Name)
	}

	created, updated, deleted, _ := splitPositions(dto.Id, dto.Positions)
	dto.Hash = CalculateHash(dto.Positions)
	oldOrder := &models.Order{}

	logger.Debug("update",
		logger.StringAttr("hash", dto.Hash),
		logger.IntAttr("len", len(dto.Positions)),
	)

	err := s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		// Получаем старое состояние заказа для логирования
		var err error
		oldOrder, err = s.GetById(ctx, tx, &models.GetOrderByIdDTO{Id: dto.Id})
		if err != nil {
			return fmt.Errorf("failed to get old order: %w", err)
		}

		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update order. error: %w", err)
		}

		if err := s.positions.Create(ctx, tx, created); err != nil {
			return err
		}
		if err := s.positions.Update(ctx, tx, updated); err != nil {
			return err
		}
		if err := s.positions.Delete(ctx, tx, deleted); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	go s.activity.AsyncLog(context.Background(), func() error {
		return s.txManager.WithinTransaction(context.Background(), func(tx postgres.Tx) error {
			if err := s.activity.LogOrderUpdate(context.Background(), dto.Actor, oldOrder, dto); err != nil {
				return fmt.Errorf("failed to log order update: %w", err)
			}
			return s.activity.BatchLogPositions(context.Background(), tx, &models.BatchLogPositionsDTO{
				OrderID: dto.Id,
				Actor:   dto.Actor,
				Created: created,
				Updated: updated,
				Deleted: deleted,
				Old:     oldOrder.Positions,
			})
		})
	}, map[string]any{
		"order_id": dto.Id,
		"actor":    dto.Actor,
		"action":   "order_" + models.ActionUpdate,
	})

	return nil
}
