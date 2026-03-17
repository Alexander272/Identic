package models

import "time"

type GetOrderByIdDTO struct {
	Id string `json:"id" db:"id"`
}

type GetOrderByYearDTO struct {
	Year int `json:"year" db:"year"`
}

type GetUniqueDTO struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}

type Order struct {
	Id            string      `json:"id" db:"id"`
	Customer      string      `json:"customer" db:"customer"`
	Consumer      string      `json:"consumer" db:"consumer"`
	Manager       string      `json:"manager" db:"manager"`
	Bill          string      `json:"bill" db:"bill"`
	Date          time.Time   `json:"date" db:"date"`
	Notes         string      `json:"notes" db:"notes"`
	PositionCount int         `json:"positionCount" db:"position_count"`
	CreatedAt     time.Time   `json:"createdAt" db:"created_at"`
	Positions     []*Position `json:"positions"`
}

type OrderDTO struct {
	Id        string         `json:"id" db:"id"`
	Customer  string         `json:"customer" db:"customer"`
	Consumer  string         `json:"consumer" db:"consumer"`
	Manager   string         `json:"manager" db:"manager"`
	Bill      string         `json:"bill" db:"bill"`
	Date      time.Time      `json:"date" db:"date"`
	Year      int            `json:"year" db:"year"`
	Notes     string         `json:"notes" db:"notes"`
	Positions []*PositionDTO `json:"positions"`
}

type DeleteOrderDTO struct {
	Id string `json:"id" db:"id"`
}
