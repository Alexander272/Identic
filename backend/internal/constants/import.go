package constants

import "github.com/Alexander272/Identic/backend/internal/models"

var ImportTemplate = &models.ImportTemplate{
	DateColumn:     1,
	CustomerColumn: 2,
	ConsumerColumn: 3,
	NameColumn:     4,
	QuantityColumn: 5,
	ManagerColumn:  6,
	BillColumn:     7,
	NotesColumn:    8,
	Count:          9,
}

const (
	MaxOrdersInBatch    = 500  // Лимит по количеству заказов
	MaxPositionsInBatch = 5000 // Лимит по количеству позиций (защита от вашего кейса)
)
