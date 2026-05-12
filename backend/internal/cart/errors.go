package cart

import "errors"

var (
	ErrCartNotFound   = errors.New("cart not found")
	ErrInDatabase     = errors.New("error in database")
	ErrItemNotFound   = errors.New("item not found")
	ErrItemNotRemoved = errors.New("item not removed")
)
