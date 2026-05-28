package api

type Page struct {
	Index int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PageQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
	// Search string `form:"search"`
	// Sort   string `form:"sort"`
}
