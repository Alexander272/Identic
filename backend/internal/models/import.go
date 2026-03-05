package models

import "mime/multipart"

type ImportDTO struct {
	File *multipart.FileHeader
}

type ImportTemplate struct {
	DateColumn     int
	CustomerColumn int
	ConsumerColumn int
	NameColumn     int
	QuantityColumn int
	ManagerColumn  int
	BillColumn     int
	NotesColumn    int
	Count          int
}
