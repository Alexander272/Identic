package services

import (
	"bytes"
	"context"
	"fmt"

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
	{Key: "current_name", Header: "Наименование СИБУР", Value: func(p *models.Price) interface{} { return p.CurrentName }},
	{Key: "new_name", Header: "Наименование СИЛУР (для проверки спецификации)", Value: func(p *models.Price) interface{} { return p.NewName }},
	{Key: "price", Header: "Цена, руб", Value: func(p *models.Price) interface{} { return p.Price }},
	{Key: "template", Header: "Шаблон", Value: func(p *models.Price) interface{} { return p.Template }},
	{Key: "note", Header: "Примечание для СИЛУР", Value: func(p *models.Price) interface{} { return p.Note }},
	{Key: "need_sibur_approval", Header: "Требуется доп.согл. с СИБУР", Value: func(p *models.Price) interface{} { return p.NeedSiburApproval }},
	{Key: "code_new_name", Header: "Код / Наименование СИЛУР (для счета в ERP)", Value: func(p *models.Price) interface{} { return p.Code + " / " + p.NewName }},
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
		if !colsMap["code_new_name"] {
			columns = append(columns, allColumns[len(allColumns)-1])
		}
	}

	sheet := "Прайс"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return nil, fmt.Errorf("failed to set sheet name: %w", err)
	}

	colWidths := map[string]float64{
		"code":                12,
		"current_name":        80,
		"new_name":            80,
		"price":               12,
		"template":            40,
		"note":                40,
		"need_sibur_approval": 20,
		"code_new_name":       50,
	}
	for i, col := range columns {
		if w, ok := colWidths[col.Key]; ok {
			colName, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				return nil, fmt.Errorf("failed to get column name: %w", err)
			}
			if err := f.SetColWidth(sheet, colName, colName, w); err != nil {
				return nil, fmt.Errorf("failed to set column width: %w", err)
			}
		}
	}

	borderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create border style: %w", err)
	}

	for i, col := range columns {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to get cell name: %w", err)
		}
		if err := f.SetCellValue(sheet, cell, col.Header); err != nil {
			return nil, fmt.Errorf("failed to set cell value: %w", err)
		}
	}

	for i, p := range positions {
		for j, col := range columns {
			cell, err := excelize.CoordinatesToCellName(j+1, i+2)
			if err != nil {
				return nil, fmt.Errorf("failed to get cell name: %w", err)
			}
			if err := f.SetCellValue(sheet, cell, col.Value(p)); err != nil {
				return nil, fmt.Errorf("failed to set cell value: %w", err)
			}
		}
	}

	if len(columns) > 0 {
		hCell, err := excelize.CoordinatesToCellName(1, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to get cell name: %w", err)
		}
		vCell, err := excelize.CoordinatesToCellName(len(columns), len(positions)+1)
		if err != nil {
			return nil, fmt.Errorf("failed to get cell name: %w", err)
		}
		if err := f.SetCellStyle(sheet, hCell, vCell, borderStyle); err != nil {
			return nil, fmt.Errorf("failed to set cell style: %w", err)
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write xlsx: %w", err)
	}
	return buf.Bytes(), nil
}
