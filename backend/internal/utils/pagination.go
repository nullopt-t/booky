package utils

type PageResult[T any] struct {
	Items []T `json:"data"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PaginationQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}
