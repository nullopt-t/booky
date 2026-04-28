package db

type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}