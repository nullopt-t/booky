package inventory

import "errors"

var (
	ErrInDatabase           = errors.New("database error")
	ErrInsufficientQuantity = errors.New("insufficient product quantity")
)
