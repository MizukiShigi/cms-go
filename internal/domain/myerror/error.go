package myerror

import (
	"errors"
	"net/http"
)

type Code string

const (
	InvalidCode             Code = "INVALID"
	InternalServerErrorCode Code = "INTERNAL_SERVER_ERROR"
	UnauthorizedCode        Code = "UNAUTHORIZED"
	ForbiddenCode           Code = "FORBIDDEN"
	NotFoundCode            Code = "NOT_FOUND"
	ConflictCode            Code = "CONFLICT"
)

var (
	InvalidRequestError = NewMyError(InvalidCode, "Invalid request")
	InternalServerError = NewMyError(InternalServerErrorCode, "Internal server error")
	ForbiddenError      = NewMyError(ForbiddenCode, "Forbidden")
	NotFoundError       = NewMyError(NotFoundCode, "Not found")
	ConflictError       = NewMyError(ConflictCode, "Conflict")
)

type MyError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func NewMyError(code Code, message string) *MyError {
	return &MyError{
		Code:    code,
		Message: message,
	}
}

func (e *MyError) Error() string {
	return e.Message
}

func (e *MyError) StatusCode() int {
	switch e.Code {
	case InvalidCode:
		return http.StatusBadRequest
	case InternalServerErrorCode:
		return http.StatusInternalServerError
	case UnauthorizedCode:
		return http.StatusUnauthorized
	case ForbiddenCode:
		return http.StatusForbidden
	case NotFoundCode:
		return http.StatusNotFound
	case ConflictCode:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func IsMyError(err error) bool {
	var domainErr *MyError
	return errors.As(err, &domainErr)
}
