package errors

import (
	// "errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("error: path or body: (%s), status code %d", e.Message, e.Code)
}

func New(code int, message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    code,
	}
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, message)
}

func UnexpectedError(statusCode int, message string) *HTTPError {
	return New(statusCode, message)
}
