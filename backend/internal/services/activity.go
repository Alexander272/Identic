package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/google/uuid"
)

type ActivityService struct {
	repo repository.Activity
}

func NewActivityService(repo repository.Activity) *ActivityService {
	return &ActivityService{
		repo: repo,
	}
}

type Activity interface {
	Get(ctx context.Context, req *models.GetAllActivityLogsDTO) ([]*models.ActivityLog, error)
	GetByOrder(ctx context.Context, orderID string) ([]*models.ActivityLog, error)
	GetByEntity(ctx context.Context, req *models.GetActivityLogsDTO) ([]*models.ActivityLog, error)
	LogOrderCreate(ctx context.Context, tx postgres.Tx, order *models.OrderDTO) error
	LogOrderUpdate(ctx context.Context, actor models.Actor, oldOrder *models.Order, newOrder *models.OrderDTO) error
	LogOrderDelete(ctx context.Context, actor models.Actor, order *models.Order) error

	BatchLogPositions(ctx context.Context, tx postgres.Tx, req *models.BatchLogPositionsDTO) error
	AsyncLog(ctx context.Context, fn func() error, debugInfo map[string]any)
}

func (s *ActivityService) Get(ctx context.Context, req *models.GetAllActivityLogsDTO) ([]*models.ActivityLog, error) {
	logs, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity logs: %w", err)
	}
	return logs, nil
}

func (s *ActivityService) GetByOrder(ctx context.Context, orderID string) ([]*models.ActivityLog, error) {
	logs, err := s.repo.GetByOrder(ctx, uuid.MustParse(orderID))
	if err != nil {
		return nil, fmt.Errorf("failed to get activity logs: %w", err)
	}
	return logs, nil
}

func (s *ActivityService) GetByEntity(ctx context.Context, req *models.GetActivityLogsDTO) ([]*models.ActivityLog, error) {
	logs, err := s.repo.GetByEntity(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity logs: %w", err)
	}
	return logs, nil
}

func (s *ActivityService) LogOrderCreate(ctx context.Context, tx postgres.Tx, order *models.OrderDTO) error {
	values := make(map[string]interface{})
	values["customer"] = order.Customer
	values["consumer"] = order.Consumer
	values["date"] = order.Date
	values["manager"] = order.Manager
	values["bill"] = order.Bill
	values["notes"] = order.Notes
	values["isBargaining"] = order.IsBargaining
	values["isBudget"] = order.IsBudget

	dto := &models.CreateActivityLogDTO{
		Action:        models.ActionInsert,
		ChangedBy:     order.Actor.ID,
		ChangedByName: order.Actor.Name,
		EntityType:    models.EntityOrder,
		EntityID:      order.Id,
		Entity:        formatOrderEntity(order),
		NewValues:     values,
	}
	return s.repo.Create(ctx, nil, dto)
}

func (s *ActivityService) LogOrderUpdate(ctx context.Context, actor models.Actor, oldOrder *models.Order, newOrder *models.OrderDTO) error {
	diff := orderFieldsDiffer(oldOrder, newOrder)
	if diff == nil {
		return nil
	}

	dto := &models.CreateActivityLogDTO{
		Action:        models.ActionUpdate,
		ChangedBy:     actor.ID,
		ChangedByName: actor.Name,
		EntityType:    models.EntityOrder,
		EntityID:      newOrder.Id,
		Entity:        formatOrderEntity(newOrder),
		OldValues:     diff.OldValues,
		NewValues:     diff.NewValues,
	}
	return s.repo.Create(ctx, nil, dto)
}

func (s *ActivityService) LogOrderDelete(ctx context.Context, actor models.Actor, order *models.Order) error {
	values := make(map[string]interface{})
	values["customer"] = order.Customer
	values["consumer"] = order.Consumer
	values["date"] = order.Date
	values["manager"] = order.Manager
	values["bill"] = order.Bill
	values["notes"] = order.Notes
	values["isBargaining"] = order.IsBargaining
	values["isBudget"] = order.IsBudget

	dto := &models.CreateActivityLogDTO{
		Action:        models.ActionDelete,
		ChangedBy:     actor.ID,
		ChangedByName: actor.Name,
		EntityType:    models.EntityOrder,
		EntityID:      order.Id,
		Entity:        formatOrderEntityByParts(order.Customer, order.Consumer),
		OldValues:     values,
	}
	return s.repo.Create(ctx, nil, dto)
}

func (s *ActivityService) BatchLogPositions(ctx context.Context, tx postgres.Tx, req *models.BatchLogPositionsDTO) error {
	logs := make([]*models.CreateActivityLogDTO, 0, len(req.Created)+len(req.Updated)+len(req.Deleted))

	oldPosMap := make(map[string]*models.Position, len(req.Old))
	for _, p := range req.Old {
		oldPosMap[p.Id] = p
	}

	for _, pos := range req.Created {
		logs = append(logs, &models.CreateActivityLogDTO{
			Action:        models.ActionInsert,
			ChangedBy:     req.Actor.ID,
			ChangedByName: req.Actor.Name,
			EntityType:    models.EntityOrderItem,
			EntityID:      pos.Id,
			Entity:        formatPositionEntity(pos),
			ParentID:      &req.OrderID,
			NewValues:     pos,
		})
	}
	for _, pos := range req.Deleted {
		logs = append(logs, &models.CreateActivityLogDTO{
			Action:        models.ActionDelete,
			ChangedBy:     req.Actor.ID,
			ChangedByName: req.Actor.Name,
			EntityType:    models.EntityOrderItem,
			EntityID:      pos.Id,
			Entity:        formatPositionEntity(pos),
			ParentID:      &req.OrderID,
			OldValues:     pos,
		})
	}
	for _, pos := range req.Updated {
		logs = append(logs, &models.CreateActivityLogDTO{
			Action:        models.ActionUpdate,
			ChangedBy:     req.Actor.ID,
			ChangedByName: req.Actor.Name,
			EntityType:    models.EntityOrderItem,
			EntityID:      pos.Id,
			Entity:        formatPositionEntity(pos),
			ParentID:      &req.OrderID,
			OldValues:     oldPosMap[pos.Id],
			NewValues:     pos,
		})
	}

	return s.repo.CreateBatch(ctx, tx, logs)
}

func (s *ActivityService) AsyncLog(ctx context.Context, fn func() error, debugInfo map[string]any) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic in analytics logger", logger.AnyAttr("recover", r))
		}
	}()

	// 1 попытка + 1 retry
	for attempt := 0; attempt < 2; attempt++ {
		err := fn()
		if err == nil {
			return
		}

		logger.Info("analytics insert failed, retrying", logger.ErrAttr(err), logger.IntAttr("attempt", attempt+1))
		time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
	}

	logger.Error("failed to save analytics log after retries",
		logger.AnyAttr("debug", debugInfo))
	error_bot.Send(nil, "failed to save analytics log", debugInfo)
}

func orderFieldsDiffer(oldOrder *models.Order, newOrder *models.OrderDTO) *models.OrderDiff {
	var diff *models.OrderDiff

	// Хелпер для ленивой инициализации и записи
	addDiff := func(key string, oldVal, newVal interface{}) {
		if diff == nil {
			diff = &models.OrderDiff{
				OldValues: make(map[string]interface{}),
				NewValues: make(map[string]interface{}),
			}
		}
		diff.OldValues[key] = oldVal
		diff.NewValues[key] = newVal
	}

	if oldOrder.Date != newOrder.Date {
		addDiff("date", oldOrder.Date, newOrder.Date)
	}
	if oldOrder.Customer != newOrder.Customer {
		addDiff("customer", oldOrder.Customer, newOrder.Customer)
	}
	if oldOrder.Consumer != newOrder.Consumer {
		addDiff("consumer", oldOrder.Consumer, newOrder.Consumer)
	}
	if oldOrder.Manager != newOrder.Manager {
		addDiff("manager", oldOrder.Manager, newOrder.Manager)
	}
	if oldOrder.Bill != newOrder.Bill {
		addDiff("bill", oldOrder.Bill, newOrder.Bill)
	}
	if oldOrder.Notes != newOrder.Notes {
		addDiff("notes", oldOrder.Notes, newOrder.Notes)
	}
	if oldOrder.IsBargaining != newOrder.IsBargaining {
		addDiff("isBargaining", oldOrder.IsBargaining, newOrder.IsBargaining)
	}
	if oldOrder.IsBudget != newOrder.IsBudget {
		addDiff("isBudget", oldOrder.IsBudget, newOrder.IsBudget)
	}

	return diff
}

func formatOrderEntity(order *models.OrderDTO) string {
	return formatOrderEntityByParts(order.Customer, order.Consumer)
}

func formatOrderEntityByParts(customer, consumer string) string {
	if customer != "" && consumer != "" {
		return consumer + " / " + customer
	}
	if customer != "" {
		return customer
	}
	if consumer != "" {
		return consumer
	}
	return "Новый заказ"
}

func formatPositionEntity(pos *models.PositionDTO) string {
	name := pos.Name
	if pos.Quantity > 0 {
		return name + " (" + fmt.Sprintf("%.2f", pos.Quantity) + " шт.)"
	}
	return name
}
