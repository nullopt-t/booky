package api

type Response[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error map[string]string `json:"error"`
}

func Success[T any](data T) Response[T] {
	return Response[T]{Data: data}
}

func Error(code, msg string) ErrorResponse {
	return ErrorResponse{Error: map[string]string{"code": code, "message": msg}}
}
