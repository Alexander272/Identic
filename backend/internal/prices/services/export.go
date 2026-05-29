package services

import (
	"bytes"
	"context"
	"fmt"

	base_models "github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/prices/models"
	"github.com/Alexander272/Identic/backend/internal/prices/repository"
	"github.com/xuri/excelize/v2"
)

type ExportService struct {
	repo repository.Prices
}

func NewExportService(repo repository.Prices) *ExportService {
	return &ExportService{repo: repo}
}

type Export interface {
	ExportXLSX(ctx context.Context, req models.ExportPriceRequest) ([]byte, error)
}

func (s *ExportService) ExportXLSX(ctx context.Context, req models.ExportPriceRequest) ([]byte, error) {
	if len(req.Queries) == 0 && len(req.Codes) == 0 {
		return nil, base_models.ErrInvalidInput
	}

	normalizedQueries := normalizeQueries(req.Queries)
	codes := req.Codes
	if codes == nil {
		codes = []string{}
	}

	positions, err := s.repo.SearchAll(ctx, normalizedQueries, codes)
	if err != nil {
		return nil, fmt.Errorf("failed to search positions for export: %w", err)
	}

	return buildXLSX(positions, req.Columns)
}

type columnSpec struct {
	Key    string
	Header string
	Value  func(p *models.Price) interface{}
}

var allColumns = []columnSpec{
	{Key: "code", Header: "КОД", Value: func(p *models.Price) interface{} { return p.Code }},
	{Key: "current_name", Header: "Текущее наименование из Прайса (на него ориентируемся)", Value: func(p *models.Price) interface{} { return p.CurrentName }},
	{Key: "new_name", Header: "Наименование АСВНСИ новой позиции (после редактирования)", Value: func(p *models.Price) interface{} { return p.NewName }},
	{Key: "price", Header: "Цена", Value: func(p *models.Price) interface{} { return p.Price }},
	{Key: "template", Header: "Шаблон", Value: func(p *models.Price) interface{} { return p.Template }},
	{Key: "note", Header: "Примечание", Value: func(p *models.Price) interface{} { return p.Note }},
	{Key: "technique", Header: "Техника", Value: func(p *models.Price) interface{} { return p.Technique }},
	{Key: "under_drawing", Header: "Под чертеж", Value: func(p *models.Price) interface{} { return p.UnderDrawing }},
}

func buildXLSX(positions []*models.Price, cols []string) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	colsMap := make(map[string]bool, len(cols))
	for _, c := range cols {
		colsMap[c] = true
	}

	var columns []columnSpec
	if len(cols) == 0 {
		columns = allColumns
	} else {
		for _, c := range allColumns {
			if colsMap[c.Key] {
				columns = append(columns, c)
			}
		}
	}

	sheet := "Sheet1"
	for i, col := range columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, col.Header)
	}

	for i, p := range positions {
		for j, col := range columns {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, col.Value(p))
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write xlsx: %w", err)
	}
	return buf.Bytes(), nil
}
