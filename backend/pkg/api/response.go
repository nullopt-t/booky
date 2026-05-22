package api

type Response[T any] struct {
	Data T `json:"data"`
}

func Success[T any](data T) Response[T] {
	return Response[T]{Data: data}
}
