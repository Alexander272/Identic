package models

import "time"

type GetPositionsByOrderIdDTO struct {
	OrderId string `json:"orderId" db:"order_id"`
}

type GetPositionsByIds struct {
	Ids []string `json:"ids" db:"ids"`
}

type Position struct {
	Id        string    `json:"id" db:"id"`
	OrderId   string    `json:"orderId" db:"order_id"`
	RowNumber int       `json:"rowNumber" db:"row_number"`
	Name      string    `json:"name" db:"name"`
	Quantity  float32   `json:"quantity" db:"quantity"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type PositionDTO struct {
	Id        string  `json:"id" db:"id"`
	OrderId   string  `json:"orderId" db:"order_id"`
	RowNumber int     `json:"rowNumber" db:"row_number"`
	Name      string  `json:"name" db:"name"`
	Search    string  `json:"search" db:"search"`
	Quantity  float32 `json:"quantity" db:"quantity"`
	Notes     string  `json:"notes" db:"notes"`
}

type DeletePositionDTO struct {
	Id string `json:"id" db:"id"`
}
