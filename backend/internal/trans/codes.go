package trans

type ApiErrCode string

const (
	CART_NOT_FOUND      ApiErrCode = "CART_NOT_FOUND"
	CART_ALREADY_EXISTS ApiErrCode = "CART_ALREADY_EXISTS"
	INTERNAL_ERROR      ApiErrCode = "INTERNAL_ERROR"
)
