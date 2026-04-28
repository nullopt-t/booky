package product

import (
	"booky-backend/internal/domain"
	"time"
)

type ProductResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Title string `json:"title" binding:"required"`
	Price int    `json:"price" binding:"required"`
	Stock int    `json:"stock" binding:"required"`
}

type PaginatedProductsResponse struct {
	Products []domain.Product `json:"products"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
	Total    int              `json:"total"`
}
