package services

import (
	"context"
	"fmt"
	"strings"

	base_models "github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/prices/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
)

type Prices interface {
	Search(ctx context.Context, req models.SearchPriceRequest) ([]*models.Price, int, error)
	BatchSave(ctx context.Context, req models.BatchSavePricesRequest) error
}

type PricesService struct {
	repo      repository.Prices
	txManager TransactionManager
}

func NewPricesService(repo repository.Prices, txManager TransactionManager) *PricesService {
	return &PricesService{repo: repo, txManager: txManager}
}

func (s *PricesService) Search(ctx context.Context, req models.SearchPriceRequest) ([]*models.Price, int, error) {
	if len(req.Queries) == 0 && len(req.Codes) == 0 {
		return nil, 0, base_models.ErrInvalidInput
	}

	normalizedQueries := normalizeQueries(req.Queries)
	codes := req.Codes
	if codes == nil {
		codes = []string{}
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 20
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	prices, total, err := s.repo.Search(ctx, normalizedQueries, codes, req.Page, req.PerPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search prices. error: %w", err)
	}

	if len(normalizedQueries) > 0 {
		for i := range prices {
			prices[i].MatchedFields = computeMatchedFields(prices[i], normalizedQueries)
		}
	}

	return prices, total, nil
}

func (s *PricesService) BatchSave(ctx context.Context, req models.BatchSavePricesRequest) error {
	if len(req.Prices) == 0 && len(req.DeleteCodes) == 0 {
		return base_models.ErrInvalidInput
	}

	priceList := make([]*models.Price, len(req.Prices))
	for i, p := range req.Prices {
		priceList[i] = &models.Price{
			Code:         p.Code,
			CurrentName:  p.CurrentName,
			NewName:      p.NewName,
			Price:        p.Price,
			Template:     p.Template,
			Note:         p.Note,
			Technique:    p.Technique,
			UnderDrawing: p.UnderDrawing,
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
	q = strings.ToLower(q)
	q = strings.ReplaceAll(q, "х", "x")
	return q
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
