package models

import "github.com/google/uuid"

type Price struct {
	ID            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	CurrentName   string    `json:"current_name"`
	NewName       string    `json:"new_name"`
	Price         float64   `json:"price"`
	Template      string    `json:"template"`
	Note          string    `json:"note"`
	Technique     string    `json:"technique"`
	UnderDrawing  string    `json:"under_drawing"`
	MatchedFields []string  `json:"matched_fields,omitempty"`
}

type SearchPriceRequest struct {
	Query   string   `json:"query"`
	Codes   []string `json:"codes"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
}

type ExportPriceRequest struct {
	Query   string   `json:"query"`
	Codes   []string `json:"codes"`
	Columns []string `json:"columns"`
}

type UpdatePrice struct {
	Code         string  `json:"code" binding:"required"`
	CurrentName  string  `json:"current_name"`
	NewName      string  `json:"new_name"`
	Price        float64 `json:"price"`
	Template     string  `json:"template"`
	Note         string  `json:"note"`
	Technique    string  `json:"technique"`
	UnderDrawing string  `json:"under_drawing"`
}

type BatchSavePricesRequest struct {
	Prices      []UpdatePrice `json:"prices" binding:"required"`
	DeleteCodes []string      `json:"delete_codes"`
}
