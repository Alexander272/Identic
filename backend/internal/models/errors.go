package models

import "errors"

var (
	ErrNoRows        = errors.New("row not found")
	ErrAlreadyExists = errors.New("already exists")

	ErrSessionEmpty = errors.New("user session not found")

	ErrNotValid = errors.New("data is not valid")

	ErrNoData = errors.New("no data")
)
