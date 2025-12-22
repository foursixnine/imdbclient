package errors

import (
	// e "errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int
	Message string
}

func (e HTTPError) Error() string {
	return e.Message
}

// func (e HTTPError) Unwrap() error {
// 	return e.Message
// }

func New(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: fmt.Sprintf("error: path or body: (%s), status code %d", message, code),
	}
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, message)
}

func UnexpectedError(statusCode int, message string) *HTTPError {
	return New(statusCode, message)
}

func (e HTTPError) Is(target error) bool {
	if t, ok := target.(*HTTPError); ok {
		return e.Code == t.Code && e.Message == t.Message // Or compare more fields
	}
	return false
}
