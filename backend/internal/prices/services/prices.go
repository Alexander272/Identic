package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	base_models "github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/prices/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
)

var reMultiSpace = regexp.MustCompile(`\s+`)

type Prices interface {
	GetAll(ctx context.Context, page, perPage int) ([]*models.Price, int, error)
	Search(ctx context.Context, req *models.SearchPriceRequest) ([]*models.Price, int, error)
	SearchAll(ctx context.Context, req *models.SearchPriceRequest) ([]*models.Price, error)
	BatchSave(ctx context.Context, req *models.BatchSavePricesRequest) error
}

type PricesService struct {
	repo      repository.Prices
	txManager TransactionManager
	searchLog PriceSearchLogs
}

func NewPricesService(repo repository.Prices, txManager TransactionManager, searchLog PriceSearchLogs) *PricesService {
	return &PricesService{repo: repo, txManager: txManager, searchLog: searchLog}
}

func (s *PricesService) GetAll(ctx context.Context, page, perPage int) ([]*models.Price, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	return s.repo.GetAll(ctx, page, perPage)
}

func (s *PricesService) Search(ctx context.Context, req *models.SearchPriceRequest) ([]*models.Price, int, error) {
	if len(req.Queries) == 0 && len(req.Codes) == 0 {
		return nil, 0, base_models.ErrInvalidInput
	}

	start := time.Now()

	normalizedQueries := normalizeQueries(req.Queries)
	codes := req.Codes
	if codes == nil {
		codes = []string{}
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 500
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	prices, total, err := s.repo.Search(ctx, normalizedQueries, codes, req.Fields, req.Page, req.PerPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search prices. error: %w", err)
	}

	if len(normalizedQueries) > 0 {
		if len(req.Fields) > 0 {
			for i := range prices {
				prices[i].MatchedFields = req.Fields
			}
		} else {
			for i := range prices {
				prices[i].MatchedFields = computeMatchedFields(prices[i], normalizedQueries)
			}
		}
	}

	if s.searchLog != nil {
		s.searchLog.LogAsync(codes, req.Queries, req.ActorID, req.ActorName, time.Since(start), total)
	}

	return prices, total, nil
}

func (s *PricesService) SearchAll(ctx context.Context, req *models.SearchPriceRequest) ([]*models.Price, error) {
	prices, err := s.repo.SearchAll(ctx, req.Queries, req.Codes)
	if err != nil {
		return nil, fmt.Errorf("failed to search all prices. error: %w", err)
	}

	return prices, nil
}

func (s *PricesService) BatchSave(ctx context.Context, req *models.BatchSavePricesRequest) error {
	if len(req.Prices) == 0 && len(req.DeleteCodes) == 0 {
		return base_models.ErrInvalidInput
	}

	priceList := make([]*models.Price, len(req.Prices))
	for i, p := range req.Prices {
		priceList[i] = &models.Price{
			Code:            p.Code,
			CurrentName:     p.CurrentName,
			NewName:         p.NewName,
			Price:           p.Price,
			Template:        p.Template,
			Note:            p.Note,
			Technique:       p.Technique,
			UnderDrawing:    p.UnderDrawing,
			SearchText:      buildSearchText(p.CurrentName, p.NewName, p.Template),
			CurrentNameNorm: normalizeQuery(p.CurrentName),
			NewNameNorm:     normalizeQuery(p.NewName),
			TemplateNorm:    normalizeQuery(p.Template),
		}
	}

	return s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		if len(priceList) > 0 {
			if err := s.repo.UpsertSeveral(ctx, tx, priceList); err != nil {
				return fmt.Errorf("failed to upsert prices. error: %w", err)
			}
		}
		if len(req.DeleteCodes) > 0 {
			if err := s.repo.DeleteSeveral(ctx, tx, req.DeleteCodes); err != nil {
				return fmt.Errorf("failed to delete prices. error: %w", err)
			}
		}
		return nil
	})
}

func normalizeQuery(q string) string {
	if q == "" {
		return ""
	}
	q = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, q)
	q = strings.ToLower(q)
	q = strings.ReplaceAll(q, "х", "x")
	q = reMultiSpace.ReplaceAllString(q, " ")
	return strings.TrimSpace(q)
}

func normalizeQueries(queries []string) []string {
	if len(queries) == 0 {
		return nil
	}
	normalized := make([]string, len(queries))
	for i, q := range queries {
		normalized[i] = normalizeQuery(q)
	}
	return normalized
}

func buildSearchText(currentName, newName, template string) string {
	parts := []string{currentName, template}
	if currentName != newName {
		parts = append(parts, newName)
	}
	text := strings.Join(parts, " ")
	return normalizeQuery(text)
}

func computeMatchedFields(p *models.Price, queries []string) []string {
	if len(queries) == 0 {
		return nil
	}

	fieldsMap := make(map[string]struct{})
	for _, query := range queries {
		if strings.Contains(normalizeQuery(p.CurrentName), query) {
			fieldsMap["current_name"] = struct{}{}
		}
		if strings.Contains(normalizeQuery(p.NewName), query) {
			fieldsMap["new_name"] = struct{}{}
		}
		if strings.Contains(normalizeQuery(p.Template), query) {
			fieldsMap["template"] = struct{}{}
		}
	}

	fields := make([]string, 0, len(fieldsMap))
	for f := range fieldsMap {
		fields = append(fields, f)
	}
	return fields
}
