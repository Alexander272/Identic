package services

import (
	"context"
	"testing"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock types
type mockOrdersRepo struct{ mock.Mock }

func (m *mockOrdersRepo) Get(ctx context.Context, req *models.OrderFilterDTO) ([]*models.Order, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.Order), args.Error(1)
}
func (m *mockOrdersRepo) GetById(ctx context.Context, tx postgres.Tx, req *models.GetOrderByIdDTO) (*models.Order, error) {
	args := m.Called(ctx, tx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}
func (m *mockOrdersRepo) GetInfoById(ctx context.Context, req *models.GetOrderByIdDTO) (*models.Order, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}
func (m *mockOrdersRepo) GetByYear(ctx context.Context, req *models.GetOrderByYearDTO) ([]*models.Order, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.Order), args.Error(1)
}
func (m *mockOrdersRepo) GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]string), args.Error(1)
}
func (m *mockOrdersRepo) GetFlatData(ctx context.Context, req *models.GetFlatOrderDTO) (*models.FlatOrderRes, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FlatOrderRes), args.Error(1)
}
func (m *mockOrdersRepo) IsExist(ctx context.Context, tx postgres.Tx, dto *models.OrderDTO) (bool, error) {
	args := m.Called(ctx, tx, dto)
	return args.Bool(0), args.Error(1)
}
func (m *mockOrdersRepo) IsExistByPos(ctx context.Context, tx postgres.Tx, dto *models.OrderDTO) (string, error) {
	args := m.Called(ctx, tx, dto)
	return args.String(0), args.Error(1)
}
func (m *mockOrdersRepo) Create(ctx context.Context, tx postgres.Tx, dto *models.OrderDTO) error {
	args := m.Called(ctx, tx, dto)
	return args.Error(0)
}
func (m *mockOrdersRepo) CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.OrderDTO) error {
	args := m.Called(ctx, tx, dto)
	return args.Error(0)
}
func (m *mockOrdersRepo) Update(ctx context.Context, tx postgres.Tx, dto *models.OrderDTO) error {
	args := m.Called(ctx, tx, dto)
	return args.Error(0)
}
func (m *mockOrdersRepo) Delete(ctx context.Context, tx postgres.Tx, dto *models.DeleteOrderDTO) error {
	args := m.Called(ctx, tx, dto)
	return args.Error(0)
}

type mockTxManager struct{ mock.Mock }

func (m *mockTxManager) WithinTransaction(ctx context.Context, fn func(tx postgres.Tx) error) error {
	args := m.Called(ctx, fn)
	if fn != nil && args.Bool(0) {
		return fn(nil)
	}
	return args.Error(1)
}

type mockPositions struct{ mock.Mock }

func (m *mockPositions) GetByOrder(ctx context.Context, tx postgres.Tx, req *models.GetPositionsByOrderIdDTO) ([]*models.Position, error) {
	args := m.Called(ctx, tx, req)
	return args.Get(0).([]*models.Position), args.Error(1)
}
func (m *mockPositions) GetByIds(ctx context.Context, req *models.GetPositionsByIds) ([]*models.Position, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.Position), args.Error(1)
}
func (m *mockPositions) Create(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error {
	args := m.Called(ctx, tx, dto)
	return args.Error(0)
}
func (m *mockPositions) Update(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error {
	args := m.Called(ctx, tx, dto)
	return args.Error(0)
}
func (m *mockPositions) Delete(ctx context.Context, tx postgres.Tx, dto []*models.PositionDTO) error {
	args := m.Called(ctx, tx, dto)
	return args.Error(0)
}
func (m *mockPositions) DeleteByOrder(ctx context.Context, tx postgres.Tx, dto *models.DeletePositionsByOrderIdDTO) error {
	args := m.Called(ctx, tx, dto)
	return args.Error(0)
}

type mockSearch struct{ mock.Mock }

func (m *mockSearch) Search(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.OrderMatchResult), args.Error(1)
}
func (m *mockSearch) SearchAndGroup(ctx context.Context, req *models.SearchRequest) ([]*models.OrderMatchResult, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.OrderMatchResult), args.Error(1)
}
func (m *mockSearch) GetCache(ctx context.Context, req *models.GetCacheDTO) ([]string, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]string), args.Error(1)
}

type mockActivity struct {
	mock.Mock
	asyncLogCalls []struct {
		ctx      context.Context
		fn       func() error
		debugInfo map[string]any
	}
}

func (m *mockActivity) Get(ctx context.Context, req *models.GetAllActivityLogsDTO) ([]*models.ActivityLog, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.ActivityLog), args.Error(1)
}
func (m *mockActivity) GetByOrder(ctx context.Context, orderID string) ([]*models.ActivityLog, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).([]*models.ActivityLog), args.Error(1)
}
func (m *mockActivity) GetByEntity(ctx context.Context, req *models.GetActivityLogsDTO) ([]*models.ActivityLog, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.ActivityLog), args.Error(1)
}
func (m *mockActivity) LogOrderCreate(ctx context.Context, tx postgres.Tx, order *models.OrderDTO) error {
	args := m.Called(ctx, tx, order)
	return args.Error(0)
}
func (m *mockActivity) LogOrderUpdate(ctx context.Context, actor models.Actor, oldOrder *models.Order, newOrder *models.OrderDTO) error {
	args := m.Called(ctx, actor, oldOrder, newOrder)
	return args.Error(0)
}
func (m *mockActivity) LogOrderDelete(ctx context.Context, actor models.Actor, order *models.Order) error {
	args := m.Called(ctx, actor, order)
	return args.Error(0)
}
func (m *mockActivity) BatchLogPositions(ctx context.Context, tx postgres.Tx, req *models.BatchLogPositionsDTO) error {
	args := m.Called(ctx, tx, req)
	return args.Error(0)
}
func (m *mockActivity) AsyncLog(ctx context.Context, fn func() error, debugInfo map[string]any) {
	// Записываем вызов, но не вызываем fn(), так как это асинхронная операция
	m.asyncLogCalls = append(m.asyncLogCalls, struct {
		ctx      context.Context
		fn       func() error
		debugInfo map[string]any
	}{ctx, fn, debugInfo})
	// Не вызываем m.Called(), чтобы избежать паники при вызове из горутины
}

func newTestOrdersService() (*OrdersService, *mockOrdersRepo, *mockTxManager, *mockPositions, *mockSearch, *mockActivity) {
	repo := new(mockOrdersRepo)
	txMgr := new(mockTxManager)
	pos := new(mockPositions)
	search := new(mockSearch)
	act := new(mockActivity)
	svc := NewOrdersService(repo, txMgr, pos, search, act)
	return svc, repo, txMgr, pos, search, act
}

func TestOrdersService_Create_Success(t *testing.T) {
	svc, repo, txMgr, pos, _, _ := newTestOrdersService()

	dto := &models.OrderDTO{
		Id:       uuid.NewString(),
		Customer: "Test Customer",
		Positions: []*models.PositionDTO{
			{Name: "Item 1", Quantity: 10},
		},
	}

	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil)
	repo.On("IsExistByPos", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
	repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	id, err := svc.Create(context.Background(), dto)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Equal(t, dto.Id, id)
	assert.NotEmpty(t, dto.Hash)
	assert.Equal(t, dto.Id, dto.Positions[0].OrderId)
}

func TestOrdersService_Create_AlreadyExists(t *testing.T) {
	svc, repo, txMgr, _, _, _ := newTestOrdersService()

	dto := &models.OrderDTO{
		Id:       uuid.NewString(),
		Customer: "Test Customer",
		Positions: []*models.PositionDTO{
			{Name: "Item 1", Quantity: 10},
		},
	}

	// Заказ уже существует - IsExistByPos возвращает ID существующего заказа
	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil)
	repo.On("IsExistByPos", mock.Anything, mock.Anything, mock.Anything).Return(dto.Id, nil)

	id, err := svc.Create(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, dto.Id, id)
}

func TestOrdersService_Create_EmptyPositions(t *testing.T) {
	svc, repo, txMgr, pos, _, _ := newTestOrdersService()

	dto := &models.OrderDTO{
		Customer: "Test Customer",
	}

	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil)
	repo.On("IsExistByPos", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
	repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	id, err := svc.Create(context.Background(), dto)
	assert.NoError(t, err)
	assert.Empty(t, dto.Hash)
	assert.NotEmpty(t, id)
}

func TestOrdersService_Create_NewOrderNoId(t *testing.T) {
	svc, repo, txMgr, pos, _, _ := newTestOrdersService()

	dto := &models.OrderDTO{
		Customer: "Test Customer",
		Positions: []*models.PositionDTO{
			{Name: "Item 1", Quantity: 10},
		},
	}

	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil)
	repo.On("IsExistByPos", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
	repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	id, err := svc.Create(context.Background(), dto)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	// ID должен быть сгенерирован, если не был передан
	assert.NotEmpty(t, dto.Id)
}

func TestOrdersService_Update_Success(t *testing.T) {
	svc, repo, txMgr, pos, _, _ := newTestOrdersService()

	orderID := uuid.NewString()
	oldOrder := &models.Order{
		Id:       orderID,
		Customer: "Old Customer",
		Positions: []*models.Position{
			{Id: "pos1", Name: "Old Item", Quantity: 5},
		},
	}

	dto := &models.OrderDTO{
		Id:       orderID,
		Customer: "New Customer",
		Positions: []*models.PositionDTO{
			{Id: "pos1", Name: "Old Item", Quantity: 5, Status: models.PositionUpdated},
			{Name: "New Item", Quantity: 3, Status: models.PositionCreated},
		},
	}

	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil)
	repo.On("GetById", mock.Anything, mock.Anything, mock.Anything).Return(oldOrder, nil)
	pos.On("GetByOrder", mock.Anything, mock.Anything, mock.Anything).Return(oldOrder.Positions, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := svc.Update(context.Background(), dto)
	assert.NoError(t, err)
	assert.NotEmpty(t, dto.Hash)
}

func TestOrdersService_Update_OrderNotFound(t *testing.T) {
	svc, repo, txMgr, _, _, _ := newTestOrdersService()

	dto := &models.OrderDTO{
		Id:       uuid.NewString(),
		Customer: "Test",
		Positions: []*models.PositionDTO{
			{Name: "Item", Quantity: 1, Status: models.PositionUpdated},
		},
	}

	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil)
	repo.On("GetById", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

	err := svc.Update(context.Background(), dto)
	assert.Error(t, err)
}

func TestOrdersService_Update_WithDeletedPositions(t *testing.T) {
	svc, repo, txMgr, pos, _, _ := newTestOrdersService()

	orderID := uuid.NewString()
	oldOrder := &models.Order{
		Id:        orderID,
		Positions: []*models.Position{{Id: "pos1", Name: "Item 1", Quantity: 5}},
	}

	dto := &models.OrderDTO{
		Id:       orderID,
		Customer: "Test",
		Positions: []*models.PositionDTO{
			{Id: "pos1", Status: models.PositionDeleted},
		},
	}

	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil)
	repo.On("GetById", mock.Anything, mock.Anything, mock.Anything).Return(oldOrder, nil)
	pos.On("GetByOrder", mock.Anything, mock.Anything, mock.Anything).Return(oldOrder.Positions, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	pos.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := svc.Update(context.Background(), dto)
	assert.NoError(t, err)
}

func TestOrdersService_Create_SameContentDifferentCustomer(t *testing.T) {
	svc, repo, txMgr, pos, _, _ := newTestOrdersService()

	positions := []*models.PositionDTO{
		{Name: "Item 1", Quantity: 10},
	}

	dto1 := &models.OrderDTO{
		Id:        uuid.NewString(),
		Customer:  "Customer A",
		Positions: positions,
	}

	dto2 := &models.OrderDTO{
		Id:        uuid.NewString(),
		Customer:  "Customer B",
		Positions: positions,
	}

	// Первый заказ создается успешно
	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil).Once()
	repo.On("IsExistByPos", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	pos.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	id1, err := svc.Create(context.Background(), dto1)
	assert.NoError(t, err)
	assert.Equal(t, dto1.Id, id1)

	// Второй заказ с тем же содержимым - возвращаем ID первого заказа
	txMgr.On("WithinTransaction", mock.Anything, mock.Anything).Return(true, nil).Once()
	repo.On("IsExistByPos", mock.Anything, mock.Anything, mock.Anything).Return(dto1.Id, nil).Once()

	id2, err := svc.Create(context.Background(), dto2)
	assert.NoError(t, err)
	assert.Equal(t, dto1.Id, id2) // Возвращаем ID существующего заказа
}

func TestOrdersService_Create_HashConsistency(t *testing.T) {
	hash1 := CalculateHash([]*models.PositionDTO{
		{Name: "Item A", Quantity: 10},
	})
	hash2 := CalculateHash([]*models.PositionDTO{
		{Name: "Item B", Quantity: 10},
	})
	hash3 := CalculateHash([]*models.PositionDTO{
		{Name: "Item A", Quantity: 20},
	})

	assert.NotEqual(t, hash1, hash2, "different names should produce different hashes")
	assert.NotEqual(t, hash1, hash3, "different quantities should produce different hashes")
}

func TestOrdersService_Create_SamePositionsDifferentOrder(t *testing.T) {
	positions1 := []*models.PositionDTO{
		{Name: "B", Quantity: 3},
		{Name: "A", Quantity: 2},
	}
	positions2 := []*models.PositionDTO{
		{Name: "A", Quantity: 2},
		{Name: "B", Quantity: 3},
	}

	hash1 := CalculateHash(positions1)
	hash2 := CalculateHash(positions2)

	assert.Equal(t, hash1, hash2, "same positions in different order should produce same hash")
}
