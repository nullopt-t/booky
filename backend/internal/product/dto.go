package product

import (
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

type ProductsResponse struct {
	Products []ProductResponse `json:"data"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Total    int               `json:"total"`
}

type CreateProductRequest struct {
	Title string `json:"title" binding:"required"`
	Price int    `json:"price" binding:"required"`
	Stock int    `json:"stock" binding:"required"`
}

type UpdateProductRequest struct {
	Title *string `json:"title,omitempty"`
	Price *int    `json:"price,omitempty"`
	Stock *int    `json:"stock,omitempty"`
}
