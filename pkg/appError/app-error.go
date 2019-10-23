package appError

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func NewWithError(err error, code int, message string) *AppError {
	return &AppError{err, message, code}
}
func NewWithErrorf(err error, code int, format string, a ...interface{}) *AppError {
	return &AppError{err, fmt.Sprintf(format, a...), code}
}
func New(code int, message string) *AppError {
	err := errors.New(message)
	return &AppError{err, message, code}
}
func Newf(code int, format string, a ...interface{}) *AppError {
	return New(code, fmt.Sprintf(format, a...))
}

func (appError *AppError) Error() string {
	return fmt.Sprintf("app error: code = %d desc = %s", appError.Code, appError.Message)
}
func CheckForTimeout(ctx context.Context) (err error) {
	select {
	case <-ctx.Done():
		err = NewWithError(ctx.Err(), http.StatusRequestTimeout, ctx.Err().Error())
	default:

	}
	return
}
