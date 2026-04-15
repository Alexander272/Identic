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
	IsFound   bool      `json:"isFound" db:"is_found"`
}

type PositionDTO struct {
	Id              string         `json:"id" db:"id"`
	OrderId         string         `json:"orderId" db:"order_id"`
	RowNumber       int            `json:"rowNumber" db:"row_number"`
	Name            string         `json:"name" db:"name"`
	Search          string         `json:"search" db:"search"`
	Quantity        float32        `json:"quantity" db:"quantity"`
	Notes           string         `json:"notes" db:"notes"`
	NormalizedNotes string         `json:"normalizedNotes" db:"normalized_notes"`
	Status          PositionStatus `json:"status" db:"status"`
}

type PositionStatus string

const (
	PositionCreated PositionStatus = "CREATED"
	PositionUpdated PositionStatus = "UPDATED"
	PositionDeleted PositionStatus = "DELETED"
)

type DeletePositionDTO struct {
	Id string `json:"id" db:"id"`
}

type DeletePositionsByOrderIdDTO struct {
	OrderId string `json:"orderId" db:"order_id"`
}
