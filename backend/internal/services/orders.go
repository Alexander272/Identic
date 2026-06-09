package services

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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
	ExportOrderXLSX(ctx context.Context, req *models.ExportOrderRequest) ([]byte, error)
	Create(ctx context.Context, dto *models.OrderDTO) (string, error)
	CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.OrderDTO) error
	Update(ctx context.Context, dto *models.OrderDTO) error
	Delete(ctx context.Context, dto *models.DeleteOrderDTO) error
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

func (s *OrdersService) ExportOrderXLSX(ctx context.Context, req *models.ExportOrderRequest) ([]byte, error) {
	order, err := s.repo.GetById(ctx, nil, &models.GetOrderByIdDTO{Id: req.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get order for export: %w", err)
	}

	positions, err := s.positions.GetByOrder(ctx, nil, &models.GetPositionsByOrderIdDTO{OrderId: req.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get positions for export: %w", err)
	}
	order.Positions = positions

	return buildOrderXLSX(order)
}

type orderInfoField struct {
	Label string
	Value func(o *models.Order) interface{}
}

var orderInfoFields = []orderInfoField{
	{Label: "Конечник", Value: func(o *models.Order) interface{} { return o.Consumer }},
	{Label: "Заказчик / Перекуп", Value: func(o *models.Order) interface{} { return o.Customer }},
	{Label: "Менеджер / помощник	", Value: func(o *models.Order) interface{} { return o.Manager }},
	{Label: "Счет в 1С", Value: func(o *models.Order) interface{} { return o.Bill }},
	{Label: "Дата", Value: func(o *models.Order) interface{} { return o.Date.Format("02.01.2006") }},
	{Label: "Тендер", Value: func(o *models.Order) interface{} {
		if o.IsBargaining {
			return "Да"
		}
		return "Нет"
	}},
	{Label: "Бюджет", Value: func(o *models.Order) interface{} {
		if o.IsBudget {
			return "Да"
		}
		return "Нет"
	}},
	{Label: "Примечание", Value: func(o *models.Order) interface{} { return o.Notes }},
}

func buildOrderXLSX(order *models.Order) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Позиции"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return nil, fmt.Errorf("failed to set sheet name: %w", err)
	}

	baseStyle := &excelize.Style{
		Font: &excelize.Font{Size: 9, Family: "Arial"},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	}

	borderStyle, err := f.NewStyle(baseStyle)
	if err != nil {
		return nil, fmt.Errorf("failed to create border style: %w", err)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Size: 9, Bold: true, Family: "Arial"},
		Border: baseStyle.Border,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %w", err)
	}

	labelStyle, err := f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Size: 9, Bold: true, Family: "Arial"},
		Border: baseStyle.Border,
		Fill:   excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create label style: %w", err)
	}

	centerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border:    baseStyle.Border,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create center style: %w", err)
	}

	// Column widths
	f.SetColWidth(sheet, "A", "A", 20) // Параметр / №
	f.SetColWidth(sheet, "B", "B", 60) // Значение / Наименование
	f.SetColWidth(sheet, "C", "C", 10) // Кол-во
	f.SetColWidth(sheet, "D", "D", 30) // Примечание

	infoRowCount := len(orderInfoFields)
	headerLine := infoRowCount + 2
	posStartLine := headerLine + 1

	// Write order info (rows 1..N)
	for i, field := range orderInfoFields {
		row := i + 1

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), field.Label)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), field.Value(order))
	}

	// Position headers (row N+2)
	posHeaders := []string{"№", "Наименование", "Кол-во", "Примечание"}
	for i, h := range posHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerLine)
		f.SetCellValue(sheet, cell, h)
	}

	hCell, _ := excelize.CoordinatesToCellName(1, headerLine)
	vCell, _ := excelize.CoordinatesToCellName(4, headerLine)
	f.SetCellStyle(sheet, hCell, vCell, headerStyle)

	// Write positions (rows N+3..)
	for i, pos := range order.Positions {
		row := posStartLine + i

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), pos.RowNumber)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), pos.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), pos.Quantity)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), pos.Notes)
	}

	// Base borders
	if infoRowCount > 0 {
		hCell, _ = excelize.CoordinatesToCellName(1, 1)
		vCell, _ = excelize.CoordinatesToCellName(2, infoRowCount)
		f.SetCellStyle(sheet, hCell, vCell, borderStyle)
	}

	posCount := len(order.Positions)
	if posCount > 0 {
		posEndLine := posStartLine + posCount - 1
		hCell, _ = excelize.CoordinatesToCellName(1, posStartLine)
		vCell, _ = excelize.CoordinatesToCellName(4, posEndLine)
		f.SetCellStyle(sheet, hCell, vCell, borderStyle)
	}

	// Granular styles on top of borders
	for i := range orderInfoFields {
		row := i + 1
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), labelStyle)
	}
	for i := range order.Positions {
		row := posStartLine + i
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), centerStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), centerStyle)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write xlsx: %w", err)
	}
	return buf.Bytes(), nil
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

func (s *OrdersService) Delete(ctx context.Context, dto *models.DeleteOrderDTO) error {
	oldOrder := &models.Order{}

	err := s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		var err error
		oldOrder, err = s.GetById(ctx, tx, &models.GetOrderByIdDTO{Id: dto.Id})
		if err != nil {
			return fmt.Errorf("failed to get old order: %w", err)
		}

		if err := s.repo.Delete(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to delete order. error: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	go s.activity.AsyncLog(context.Background(), func() error {
		return s.txManager.WithinTransaction(context.Background(), func(tx postgres.Tx) error {
			if err := s.activity.LogOrderDelete(context.Background(), dto.Actor, oldOrder); err != nil {
				return fmt.Errorf("failed to log order delete: %w", err)
			}
			return nil
		})
	}, map[string]any{
		"order_id": dto.Id,
		"action":   "order_" + models.ActionDelete,
		"actor":    dto.Actor,
	})

	return nil
}
