package models

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Price struct {
	ID               uuid.UUID `json:"id"`
	Code             string    `json:"code"`
	CurrentName      string    `json:"current_name"`
	NewName          string    `json:"new_name"`
	Price            float64   `json:"price"`
	Template         string    `json:"template"`
	Note             string    `json:"note"`
	NeedSiburApproval string   `json:"need_sibur_approval"`
	MatchedFields    []string  `json:"matched_fields,omitempty"`

	CurrentNameNorm string `json:"-"`
	NewNameNorm     string `json:"-"`
	TemplateNorm    string `json:"-"`
	SearchText      string `json:"-"`
}

type SearchPriceRequest struct {
	Queries   []string  `form:"query" json:"queries"`
	Codes     []string  `json:"codes"`
	Fields    []string  `json:"fields"`
	Page      int       `json:"page"`
	PerPage   int       `json:"limit"`
	ActorID   uuid.UUID `json:"-"`
	ActorName string    `json:"-"`
}

type ExportPriceRequest struct {
	Queries []string `form:"query" json:"queries"`
	Codes   []string `json:"codes"`
	Columns []string `json:"columns"`
}

type UpdatePrice struct {
	Code              string  `json:"code" binding:"required"`
	CurrentName       string  `json:"current_name"`
	NewName           string  `json:"new_name"`
	Price             float64 `json:"price"`
	Template          string  `json:"template"`
	Note              string  `json:"note"`
	NeedSiburApproval string  `json:"need_sibur_approval"`
}

type BatchSavePricesRequest struct {
	Prices      []*UpdatePrice `json:"prices" binding:"required"`
	DeleteCodes []string       `json:"delete_codes"`
}

func CleanAndValidateAtLeastOne(sl validator.StructLevel) {
	// Получаем указатель на структуру, чтобы была возможность изменить её значения
	req, ok := sl.Current().Addr().Interface().(*SearchPriceRequest)
	if !ok {
		return
	}

	// Фильтруем слайсы на месте
	req.Queries = filterEmptyStrings(req.Queries)
	req.Codes = filterEmptyStrings(req.Codes)

	// Если после очистки оба слайса оказались пустыми — возвращаем ошибку
	if len(req.Queries) == 0 && len(req.Codes) == 0 {
		sl.ReportError(req.Queries, "Queries", "queries", "at_least_one", "")
	}
}

// Вспомогательная функция для удаления пустых строк и строк из одних пробелов
func filterEmptyStrings(slice []string) []string {
	var result []string
	for _, val := range slice {
		// strings.TrimSpace уберет пробелы
		if strings.TrimSpace(val) != "" {
			result = append(result, val)
		}
	}
	return result
}
