package errors

import (
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
		return e.Code == t.Code && e.Message == t.Message
	}
	return false
}

type IMDBClientError struct {
	Err     error
	Message string
}

func (e IMDBClientError) Error() string {
	return e.Message
}

func (e IMDBClientError) Unwrap() error {
	return e.Err
}

func NewIMDBClientGenericError(message string, err error) *IMDBClientError {
	return &IMDBClientError{
		Err:     err,
		Message: fmt.Sprintf("IMDB Client Error: (%s)", message),
	}
}

type IMDBClientApplicationError struct {
	ClientError error
	AppMessage  string
}

func (e IMDBClientApplicationError) Error() string {
	return fmt.Sprintf("IMDB Client Application Error: %s", e.AppMessage)
}

func (e IMDBClientApplicationError) Unwrap() error {
	return e.ClientError
}

func NewIMDBClientApplicationError(appMessage string, clientError error) *IMDBClientApplicationError {
	return &IMDBClientApplicationError{
		ClientError: clientError,
		AppMessage:  appMessage,
	}
}
