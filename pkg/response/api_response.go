package response

// ApiResponse represents a standardized API response structure
type ApiResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
}

// Success returns a successful API response with the given data and message.
// If no message is provided, it defaults to "Success".
// It is a generic function that can be used to return successful responses for any API endpoint.
// The message parameter is variadic, allowing you to pass in any number of string arguments.
// If more than one string is passed, only the first one is used.
// The function returns an ApiResponse[T] struct, which contains the data and message.
// The Data field contains the data returned by the API endpoint, and the Message field contains the message.
// The Message field is optional, and if not provided, it defaults to "Success".
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
