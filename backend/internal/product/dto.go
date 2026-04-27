package product

type CreateProductRequest struct {
	Title string `json:"title"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}
