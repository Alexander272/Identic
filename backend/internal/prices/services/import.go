package services

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/prices/repository"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type ImportService struct {
	repo      repository.Prices
	txManager TransactionManager
}

func NewImportService(repo repository.Prices, txManager TransactionManager) *ImportService {
	return &ImportService{repo: repo, txManager: txManager}
}

type Import interface {
	ImportXLSX(ctx context.Context, r io.Reader) error
}

func (s *ImportService) ImportXLSX(ctx context.Context, r io.Reader) error {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return fmt.Errorf("failed to open xlsx file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return fmt.Errorf("failed to read sheet: %w", err)
	}

	if len(rows) < 2 {
		return nil
	}

	dataRows := rows[1:]
	data := make([]*models.Price, 0, len(dataRows))

	for _, row := range dataRows {
		if len(row) < 2 {
			continue
		}
		priceRaw := strings.TrimSpace(safeGet(row, 3))
		priceStr := strings.ReplaceAll(priceRaw, " ", "")
		if strings.Contains(priceStr, ",") && strings.Contains(priceStr, ".") {
			priceStr = strings.ReplaceAll(priceStr, ",", "")
		} else if strings.Contains(priceStr, ",") {
			priceStr = strings.ReplaceAll(priceStr, ",", ".")
		}
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			logger.Warn("failed to parse price",
				logger.StringAttr("code", strings.TrimSpace(row[1])),
				logger.StringAttr("raw", safeGet(row, 3)),
			)
			price = 0
		}

		p := &models.Price{
			Code:              strings.TrimSpace(row[0]),
			CurrentName:       safeGet(row, 1),
			NewName:           safeGet(row, 2),
			Price:             price,
			Template:          safeGet(row, 4),
			Note:              safeGet(row, 5),
			NeedSiburApproval: safeGet(row, 6),
			SearchText:        buildSearchText(safeGet(row, 1), safeGet(row, 2), safeGet(row, 4)),
			CurrentNameNorm:   normalizeQuery(safeGet(row, 1)),
			NewNameNorm:       normalizeQuery(safeGet(row, 2)),
			TemplateNorm:      normalizeQuery(safeGet(row, 4)),
		}
		data = append(data, p)
	}

	return s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		return s.repo.CreateSeveral(ctx, tx, data)
	})
}

func safeGet(row []string, idx int) string {
	if idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}
