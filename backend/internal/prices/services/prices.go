package services

import (
	"context"
	"fmt"
	"strconv"
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
	if req.Query == "" && len(req.Codes) == 0 {
		return nil, 0, base_models.ErrInvalidInput
	}

	normalizedQuery := normalizeQuery(req.Query)
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

	prices, total, err := s.repo.Search(ctx, normalizedQuery, codes, req.Page, req.PerPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search prices. error: %w", err)
	}

	if normalizedQuery != "" {
		for i := range prices {
			prices[i].MatchedFields = computeMatchedFields(prices[i], normalizedQuery)
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

func computeMatchedFields(p *models.Price, query string) []string {
	if query == "" {
		return nil
	}
	var fields []string
	if strings.Contains(normalizeQuery(p.CurrentName), query) {
		fields = append(fields, "current_name")
	}
	if strings.Contains(normalizeQuery(p.NewName), query) {
		fields = append(fields, "new_name")
	}
	if strings.Contains(normalizeQuery(strconv.FormatFloat(p.Price, 'f', -1, 64)), query) {
		fields = append(fields, "price")
	}
	if strings.Contains(normalizeQuery(p.Template), query) {
		fields = append(fields, "template")
	}
	return fields
}
