package api

type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type PageQuery struct {
	Page     int `form:"page,default:1,min:1,max:1000"`
	PageSize int `form:"pageSize,default:10,min:1,max:100"`
}
