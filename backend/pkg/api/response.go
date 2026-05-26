package api

type Response[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error map[string]string `json:"error"`
}

type PaginatedResponse[T any] struct {
	Data T   `json:"data"`
	Meta any `json:"meta"`
}

func Success[T any](data T) Response[T] {
	return Response[T]{Data: data}
}

func SuccessPaginated[T any](data T, meta any) PaginatedResponse[T] {
	return PaginatedResponse[T]{Data: data, Meta: meta}
}

func Error(code, msg string) ErrorResponse {
	return ErrorResponse{Error: map[string]string{"code": code, "message": msg}}
}
