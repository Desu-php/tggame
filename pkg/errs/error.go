package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf(e.Message)
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}
