package product

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	Total      int                `json:"total"`
}

type ProductResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Price     int       `json:"price"`
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
}

type UpdateProductRequest struct {
	Title *string `json:"title,omitempty"`
	Price *int    `json:"price,omitempty"`
}
