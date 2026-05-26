package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID
	Title     string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductCategory struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
