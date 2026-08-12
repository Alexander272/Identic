package services

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockSearchRepo имитирует repository.Search (postgres + redis части)
type mockSearchRepo struct{ mock.Mock }

func (m *mockSearchRepo) FetchExact(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.RawMatch), args.Error(1)
}
func (m *mockSearchRepo) FetchFuzzy(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.RawMatch), args.Error(1)
}
func (m *mockSearchRepo) FetchExactByQuantity(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.RawMatch), args.Error(1)
}
func (m *mockSearchRepo) FetchFuzzyByQuantity(ctx context.Context, req *models.SearchRequest) ([]*models.RawMatch, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*models.RawMatch), args.Error(1)
}
func (m *mockSearchRepo) GetCache(ctx context.Context, req *models.GetCacheDTO) ([]string, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]string), args.Error(1)
}
func (m *mockSearchRepo) SetCache(ctx context.Context, req *models.SetCacheDTO) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// searchTokenSplit повторяет regexp_split_to_array из FetchFuzzy (postgres/search.go)
var searchTokenSplit = regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9.-]+`)

func searchTokens(name string) []string {
	parts := searchTokenSplit.Split(name, -1)
	tokens := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			tokens = append(tokens, p)
		}
	}
	return tokens
}

func newTestSearchService() (*SearchService, *mockSearchRepo) {
	repo := new(mockSearchRepo)
	svc := NewSearchService(repo, "http://localhost:8080", time.Minute)
	return svc, repo
}

// Полное наименование целевой позиции (размер 3 139 121 109 4 5).
// Вводится так, как пишет пользователь: латиница °С и СМ2.
const fullNameLatin = "ПРОКЛАДКА СНП-В-3-139-121-109-4,5 ОСТ26.260.454-99+ТЕРМОРАСШИРЕННЫЙ ГРАФИТ 43,2-131КГ/СМ2 314-399°С"

// То же наименование после NormalizeString (search.go:42)
const fullNameNormalized = "прокладка снп в 3 139 121 109 4 5 ост26 260 454 99 терморасширенный графит 43 2 131кг см2 314 399 с"

func rawMatchForFullName() *models.RawMatch {
	return &models.RawMatch{
		OrderId:      "order-600",
		YearInt:      2024,
		Customer:     "Заказчик",
		Consumer:     "Потребитель",
		Date:         time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		IsBargaining: false,
		IsBudget:     false,
		ReqId:        "1",
		PosId:        "f4c80252-1c12-4325-a89a-de84cf1280a5",
		PSearch:      fullNameNormalized + " 600 шт",
		ReqTokens:    searchTokens(fullNameNormalized),
		ReqQty:       600,
		DbQty:        600,
		Similarity:   0.9,
	}
}

func TestSearchFlow_Exact_FullNameWithLatin_FindsByQty(t *testing.T) {
	svc, repo := newTestSearchService()

	req := &models.SearchRequest{
		Items:   []models.SearchItem{{Id: 1, Name: fullNameLatin, Quantity: 600}},
		IsFuzzy: false,
	}

	repo.On("FetchExact", mock.Anything, mock.MatchedBy(func(r *models.SearchRequest) bool {
		return len(r.Items) == 1 && r.Items[0].Name == fullNameNormalized
	})).Return([]*models.RawMatch{rawMatchForFullName()}, nil)
	repo.On("SetCache", mock.Anything, mock.Anything).Return(nil)

	results, err := svc.Search(context.Background(), req)
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	requireOneResult(t, results, "order-600", true)
}

func TestSearchFlow_Fuzzy_FullNameWithLatin_FindsByQty(t *testing.T) {
	svc, repo := newTestSearchService()

	req := &models.SearchRequest{
		Items:   []models.SearchItem{{Id: 1, Name: fullNameLatin, Quantity: 600}},
		IsFuzzy: true,
	}

	repo.On("FetchFuzzy", mock.Anything, mock.MatchedBy(func(r *models.SearchRequest) bool {
		return len(r.Items) == 1 && r.Items[0].Name == fullNameNormalized
	})).Return([]*models.RawMatch{rawMatchForFullName()}, nil)
	repo.On("SetCache", mock.Anything, mock.Anything).Return(nil)

	results, err := svc.Search(context.Background(), req)
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	requireOneResult(t, results, "order-600", true)
}

// Длинная фраза (размер 43 + аннотация), внутри которой есть более короткая запись.
const longPhrase = "прокладка снп в 3 43 1 6 4 5 ост 26 260 454 99 терморасширенный графит 43 2 131кг см2 314 399 с"

const shortRecordName = "прокладка снп в 3 43 1 6 4 5 ост 26 260 454 99"

// Исходная проблема: поиск по длинной фразе должен находить короткую запись, которая в неё входит.
func TestSearchFlow_Fuzzy_LongPhrase_FindsShortContainedRecord(t *testing.T) {
	svc, repo := newTestSearchService()

	req := &models.SearchRequest{
		Items:   []models.SearchItem{{Id: 1, Name: longPhrase, Quantity: 47}},
		IsFuzzy: true,
	}

	// SQL возвращает ReqTokens = токены запроса, PSearch = наименование записи в БД
	shortRecord := &models.RawMatch{
		OrderId:      "order-47",
		YearInt:      2024,
		Customer:     "Заказчик",
		Consumer:     "Потребитель",
		Date:         time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		IsBargaining: false,
		IsBudget:     false,
		ReqId:        "1",
		PosId:        "pos-47",
		PSearch:      shortRecordName,
		ReqTokens:    searchTokens(longPhrase),
		ReqQty:       47,
		DbQty:        47,
		Similarity:   0.52,
	}

	repo.On("FetchFuzzy", mock.Anything, mock.Anything).Return([]*models.RawMatch{shortRecord}, nil)
	repo.On("SetCache", mock.Anything, mock.Anything).Return(nil)

	results, err := svc.Search(context.Background(), req)
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	requireOneResult(t, results, "order-47", true)
}

// Кандидат неверного типа (д вместо в) должен отсеиваться буквой типа СНП
func TestSearchFlow_Fuzzy_RejectsWrongTypeCandidate(t *testing.T) {
	svc, repo := newTestSearchService()

	req := &models.SearchRequest{
		Items:   []models.SearchItem{{Id: 1, Name: longPhrase, Quantity: 47}},
		IsFuzzy: true,
	}

	wrongType := &models.RawMatch{
		OrderId:      "order-wrong-type",
		YearInt:      2024,
		Customer:     "Заказчик",
		Consumer:     "Потребитель",
		Date:         time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		IsBargaining: false,
		IsBudget:     false,
		ReqId:        "1",
		PosId:        "pos-wrong-type",
		PSearch:      "прокладка снп д 3 43 1 6 4 5 304л ост26 260 454 99 50 00 шт",
		ReqTokens:    searchTokens(longPhrase),
		ReqQty:       47,
		DbQty:        50,
		Similarity:   0.7,
	}
	rightType := &models.RawMatch{
		OrderId:      "order-47",
		YearInt:      2024,
		Customer:     "Заказчик",
		Consumer:     "Потребитель",
		Date:         time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		IsBargaining: false,
		IsBudget:     false,
		ReqId:        "1",
		PosId:        "pos-47",
		PSearch:      shortRecordName,
		ReqTokens:    searchTokens(longPhrase),
		ReqQty:       47,
		DbQty:        47,
		Similarity:   0.52,
	}

	repo.On("FetchFuzzy", mock.Anything, mock.Anything).Return([]*models.RawMatch{wrongType, rightType}, nil)
	repo.On("SetCache", mock.Anything, mock.Anything).Return(nil)

	results, err := svc.Search(context.Background(), req)
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	assert.Len(t, results, 1, "wrong type candidate should be filtered out")
	if len(results) > 0 {
		assert.Equal(t, "order-47", results[0].OrderId)
	}
}

func requireOneResult(t *testing.T, results []*models.OrderMatchResult, orderID string, quantEqual bool) {
	t.Helper()
	require := assert.New(t)
	if !require.Len(results, 1) {
		return
	}
	require.Equal(orderID, results[0].OrderId)
	require.Equal(1, results[0].MatchedPos)
	if require.Len(results[0].Positions, 1) {
		require.Equal(quantEqual, results[0].Positions[0].QuantEqual)
	}
	require.Equal(1, results[0].MatchedQuant)
}
