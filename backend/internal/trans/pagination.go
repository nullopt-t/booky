package trans

type Page struct {
	Index int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PaginationQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}
