package models

import (
	"time"
)

type GetOrderByIdDTO struct {
	Id          string   `json:"id" db:"id"`
	SearchId    string   `json:"searchId" db:"search_id"`
	PositionIds []string `json:"positionIds"`
}

type GetOrderByYearDTO struct {
	Year int `json:"year" db:"year"`
}

type GetUniqueDTO struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}

type GetFlatOrderDTO struct {
	Search *Search `json:"search"`
	Sort   *Sort   `json:"sort"`
	Cursor string  `json:"cursor"`
	Page   *Page   `json:"page"`
}

type Order struct {
	Id            string      `json:"id" db:"id"`
	Customer      string      `json:"customer" db:"customer"`
	Consumer      string      `json:"consumer" db:"consumer"`
	Manager       string      `json:"manager" db:"manager"`
	IsBargaining  bool        `json:"isBargaining" db:"is_bargaining"`
	IsBudget      bool        `json:"isBudget" db:"is_budget"`
	Bill          string      `json:"bill" db:"bill"`
	Date          time.Time   `json:"date" db:"date"`
	Notes         string      `json:"notes" db:"notes"`
	Year          int         `json:"year" db:"year"`
	PositionCount int         `json:"positionCount" db:"position_count"`
	CreatedAt     time.Time   `json:"createdAt" db:"created_at"`
	Positions     []*Position `json:"positions"`
	PosWereFound  bool        `json:"posWereFound"`
}

type OrderFilterDTO struct {
	Filters []*Filter `json:"filters"`
}

type OrderDTO struct {
	Id           string `json:"id" db:"id"`
	Actor        Actor
	Customer     string         `json:"customer" db:"customer"`
	Consumer     string         `json:"consumer" db:"consumer"`
	Manager      string         `json:"manager" db:"manager"`
	IsBargaining bool           `json:"isBargaining" db:"is_bargaining"`
	IsBudget     bool           `json:"isBudget" db:"is_budget"`
	Bill         string         `json:"bill" db:"bill"`
	Date         time.Time      `json:"date" db:"date"`
	Year         int            `json:"year" db:"year"`
	Notes        string         `json:"notes" db:"notes"`
	Hash         string         `json:"hash" db:"hash"`
	Positions    []*PositionDTO `json:"positions"`
}

type DeleteOrderDTO struct {
	Id string `json:"id" db:"id"`
}

type OrderUpdateEvent struct {
	Action string   `json:"action"` // "INSERT_MANY"
	Count  int      `json:"count"`
	IDs    []string `json:"ids"` // Список ID созданных заказов
}

type FlatOrder struct {
	Id            string    `json:"id" db:"id"`
	Customer      string    `json:"customer" db:"customer"`
	Consumer      string    `json:"consumer" db:"consumer"`
	Manager       string    `json:"manager" db:"manager"`
	IsBargaining  bool      `json:"isBargaining" db:"is_bargaining"`
	IsBudget      bool      `json:"isBudget" db:"is_budget"`
	Bill          string    `json:"bill" db:"bill"`
	Date          time.Time `json:"date" db:"date"`
	Notes         string    `json:"notes" db:"notes"`
	RowNumber     int       `json:"rowNumber" db:"row_number"`
	Name          string    `json:"name" db:"name"`
	Quantity      float32   `json:"quantity" db:"quantity"`
	PositionNotes string    `json:"positionNotes" db:"pos_notes"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

type FlatOrderRes struct {
	Orders  []*FlatOrder `json:"orders"`
	Cursor  string       `json:"cursor"`
	Count   int          `json:"count"`
	HasMore bool         `json:"hasMore"`
}

// CursorValue возвращает значение и тип для курсора по имени поля
// Возвращает (nil, "") если поле не найдено
func (o *FlatOrder) CursorValue(field string) (interface{}, string, bool) {
	switch field {
	case "date":
		return o.Date, "time", true
	case "customer", "consumer", "manager", "bill", "name", "notes", "pos_notes":
		// Для текстовых полей можно вернуть одно и то же
		var val string
		switch field {
		case "customer":
			val = o.Customer
		case "consumer":
			val = o.Consumer
		case "manager":
			val = o.Manager
		case "bill":
			val = o.Bill
		case "name":
			val = o.Name
		case "notes":
			val = o.Notes
		case "pos_notes":
			val = o.PositionNotes
		}
		return val, "text", true
	case "row_number":
		return o.RowNumber, "int", true
	case "quantity":
		return o.Quantity, "float32", true
	case "p.id", "id":
		return o.Id, "uuid", true
	default:
		return nil, "", false
	}
}
