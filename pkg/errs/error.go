package errs

import (
	"errors"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type APIError struct {
	Code   string
	Params *map[string]interface{}
}

func (e *APIError) Error() string {
	return e.Code
}

func NewAPIError(code string, params *map[string]interface{}) *APIError {
	return &APIError{
		Code:   code,
		Params: params,
	}
}
