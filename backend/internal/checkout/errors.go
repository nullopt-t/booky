package checkout

import "errors"

var (
	ErrEmptyCart            = errors.New("empty cart")
	ErrNotFound             = errors.New("cart doesn't exist")
	ErrProductNotFound      = errors.New("product doesn't exist")
	ErrInsufficientQuantity = errors.New("insufficient product quantity")
)
