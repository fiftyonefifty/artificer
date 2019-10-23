package appError

import (
	"fmt"
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func New(err error, message string, code int) *AppError {
	return &AppError{err, message, code}
}
func (appError *AppError) Error() string {
	return fmt.Sprintf("app error: code = %d desc = %s", appError.Code, appError.Message)
}
