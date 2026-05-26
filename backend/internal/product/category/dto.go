package category

import "github.com/google/uuid"

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
