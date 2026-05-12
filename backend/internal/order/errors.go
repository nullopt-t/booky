package order

import "errors"

var (
	ErrInDatabase             = errors.New("database error")
	ErrNoItems                = errors.New("no items in order")
	ErrInsufficientQuanity    = errors.New("insufficient quantity")
	ErrInvalidQuantity        = errors.New("invalid quantity")
	ErrProductNotFound        = errors.New("product not found")
	ErrInvalidProductID       = errors.New("invalid product ID")
	ErrOrderNotFound          = errors.New("order not found")
	ErrOrderNotPending        = errors.New("order is not pending")
	ErrOrderAlreadyCancelled  = errors.New("order already cancelled")
	ErrOrderAlreadyConfirmed  = errors.New("order already confirmed")
	ErrInvalidOrderTransition = errors.New("invalid order transition")
)
