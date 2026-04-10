package models

import "errors"

var (
	ErrNoRows             = errors.New("row not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrOrderAlreadyExists = errors.New("order already exists")

	ErrSessionEmpty = errors.New("user session not found")

	ErrNotValid = errors.New("data is not valid")

	ErrNoData = errors.New("no data")

	ErrFieldNotAllowed = errors.New("field is not allowed")

	ErrInvalidInput      = errors.New("invalid input data")
	ErrInvalidPermission = errors.New("invalid permission")

	ErrReservedRole          = errors.New("cannot create or update reserved role")
	ErrCircularInheritance   = errors.New("circular inheritance detected")
	ErrCannotInheritFromSelf = errors.New("role cannot inherit from itself")
	ErrParentRoleNotFound    = errors.New("parent role not found or inactive")
)
