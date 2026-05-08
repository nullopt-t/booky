package product

import "time"

type Product struct {
	ID        string
	Title     string
	Price     int
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
