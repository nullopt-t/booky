package inventory

import "errors"

var (
	ErrInDatabase           = errors.New("database error")
	ErrNotFound             = errors.New("not found")
	ErrInsufficientQuantity = errors.New("insufficient product quantity")
)
