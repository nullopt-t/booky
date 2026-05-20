package product

import "errors"

var (
	ErrInDatabase      = errors.New("database error")
	ErrProductNotFount = errors.New("product doesn't exist")
)
