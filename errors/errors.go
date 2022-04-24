package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status  int
	Message string
}

func (e *Error) Code() int {
	return e.Status
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(status int, message string) *Error {
	return &Error{
		status,
		message,
	}
}

func NewNotFound(message string) *Error {
	return NewError(http.StatusNotFound, message)
}

var InternalError = NewError(http.StatusInternalServerError, "Server Internal Error")

// Deprecated: use ctx.Err() instead
var TaskCanceled = fmt.Errorf("task got canceled")
