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

	positions, err := s.repo.SearchAll(ctx, normalizedQueries, codes, req.Fields)
	if err != nil {
		return nil, fmt.Errorf("failed to search positions for export: %w", err)
	}

	if len(req.Codes) > 0 {
		codeMap := make(map[string][]*models.Price)
		for _, p := range positions {
			codeMap[p.Code] = append(codeMap[p.Code], p)
		}

		ordered := make([]*models.Price, 0, len(req.Codes))
		codeIdx := make(map[string]int)
		for _, code := range req.Codes {
			idx := codeIdx[code]
			if idx < len(codeMap[code]) {
				ordered = append(ordered, codeMap[code][idx])
				codeIdx[code] = idx + 1
			} else {
				ordered = append(ordered, &models.Price{
					Code:        code,
					CurrentName: "Не найдено",
					NotFound:    true,
				})
			}
		}
		positions = ordered
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
	{Key: "currentName", Header: "Наименование СИБУР", Value: func(p *models.Price) interface{} { return p.CurrentName }},
	{Key: "newName", Header: "Наименование СИЛУР (для проверки спецификации)", Value: func(p *models.Price) interface{} { return p.NewName }},
	{Key: "price", Header: "Цена, руб", Value: func(p *models.Price) interface{} {
		if p.NotFound {
			return ""
		}
		return p.Price
	}},
	{Key: "template", Header: "Шаблон", Value: func(p *models.Price) interface{} { return p.Template }},
	{Key: "note", Header: "Примечание для СИЛУР", Value: func(p *models.Price) interface{} { return p.Note }},
	{Key: "needSiburApproval", Header: "Требуется доп.согл. с СИБУР", Value: func(p *models.Price) interface{} { return p.NeedSiburApproval }},
	{Key: "codeNewName", Header: "Код / Наименование СИЛУР (для счета в ERP)", Value: func(p *models.Price) interface{} { return p.Code + " / " + p.NewName }},
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
		"code":              12,
		"currentName":       40,
		"newName":           40,
		"price":             12,
		"template":          20,
		"note":              20,
		"needSiburApproval": 15,
		"codeNewName":       40,
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

	redStyle, err := f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Size: 9, Family: "Arial"},
		Fill:   excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FFCCCC"}},
		Border: baseStyle.Border,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create red style: %w", err)
	}

	centerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border:    baseStyle.Border,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create center style: %w", err)
	}

	redCenterStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Family: "Arial"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FFCCCC"}},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border:    baseStyle.Border,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create red center style: %w", err)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Size: 9, Bold: true, Family: "Arial"},
		Border: baseStyle.Border,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %w", err)
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

	codeCol, priceCol := -1, -1
	for j, col := range columns {
		switch col.Key {
		case "code":
			codeCol = j + 1
		case "price":
			priceCol = j + 1
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

		for i := 2; i <= len(positions)+1; i++ {
			if codeCol > 0 {
				cell, _ := excelize.CoordinatesToCellName(codeCol, i)
				f.SetCellStyle(sheet, cell, cell, centerStyle)
			}
			if priceCol > 0 {
				cell, _ := excelize.CoordinatesToCellName(priceCol, i)
				f.SetCellStyle(sheet, cell, cell, centerStyle)
			}
		}

		for i, p := range positions {
			if p.NotFound {
				hCell, err := excelize.CoordinatesToCellName(1, i+2)
				if err != nil {
					return nil, fmt.Errorf("failed to get cell name: %w", err)
				}
				vCell, err := excelize.CoordinatesToCellName(len(columns), i+2)
				if err != nil {
					return nil, fmt.Errorf("failed to get cell name: %w", err)
				}
				if err := f.SetCellStyle(sheet, hCell, vCell, redStyle); err != nil {
					return nil, fmt.Errorf("failed to set red style: %w", err)
				}

				if codeCol > 0 {
					cell, _ := excelize.CoordinatesToCellName(codeCol, i+2)
					f.SetCellStyle(sheet, cell, cell, redCenterStyle)
				}
				if priceCol > 0 {
					cell, _ := excelize.CoordinatesToCellName(priceCol, i+2)
					f.SetCellStyle(sheet, cell, cell, redCenterStyle)
				}
			}
		}
		if len(columns) > 0 {
			hCell, _ := excelize.CoordinatesToCellName(1, 1)
			vCell, _ := excelize.CoordinatesToCellName(len(columns), 1)
			f.SetCellStyle(sheet, hCell, vCell, headerStyle)

			vCell, _ = excelize.CoordinatesToCellName(len(columns), len(positions)+1)
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write xlsx: %w", err)
	}
	return buf.Bytes(), nil
}
