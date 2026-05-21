package cart

import "errors"

var (
	ErrCartNotFound     = errors.New("cart not found")
	ErrDatabaseFailure  = errors.New("database operation failed")
	ErrDatabaseTimeout  = errors.New("database timeout")
	ErrItemNotFound     = errors.New("item not found")
	ErrItemNotRemoved   = errors.New("item not removed")
	ErrCartAlreadyExist = errors.New("cart already exists")
)
