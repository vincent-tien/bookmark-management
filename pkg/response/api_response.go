package response

type ApiResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
}

func Success[T any](data T, message ...string) ApiResponse[T] {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}
	return ApiResponse[T]{
		Data:    data,
		Message: msg,
	}
}
