package apperror

import (
	"encoding/json"
	"fmt"
)

var (
	ErrNotFound = NewAppError("not found", "AE-000404")
)

type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func NewAppError(message, code string) *AppError {
	return &AppError{
		Err:     fmt.Errorf(message),
		Code:    code,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func BadRequestError(message string) *AppError {
	return NewAppError(message, "AE-000500")
}
