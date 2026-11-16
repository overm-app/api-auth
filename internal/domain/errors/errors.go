package errors

import "net/http"

type AppError struct {
	Code       ErrorCode
	Message    string
	HTTPStatus int
	Err        error
	Details    map[string]any
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) WithDetails(d map[string]any) *AppError {
	cp := *e
	cp.Details = d
	return &cp
}

func New(code ErrorCode, status int, msg string, cause error) *AppError {
	return &AppError{Code: code, HTTPStatus: status, Message: msg, Err: cause}
}

func Validation(code ErrorCode, msg string) *AppError {
	return New(code, http.StatusBadRequest, msg, nil)
}
func Unauthorized(code ErrorCode, msg string) *AppError {
	return New(code, http.StatusUnauthorized, msg, nil)
}
func Forbidden(code ErrorCode, msg string) *AppError {
	return New(code, http.StatusForbidden, msg, nil)
}
func NotFound(code ErrorCode, msg string) *AppError {
	return New(code, http.StatusNotFound, msg, nil)
}
func Conflict(code ErrorCode, msg string) *AppError {
	return New(code, http.StatusConflict, msg, nil)
}
func Unprocessable(code ErrorCode, msg string) *AppError {
	return New(code, http.StatusUnprocessableEntity, msg, nil)
}
func BadGateway(code ErrorCode, msg string, cause error) *AppError {
	return New(code, http.StatusBadGateway, msg, cause)
}
func Internal(msg string, cause error) *AppError {
	return New(ErrInternal, http.StatusInternalServerError, msg, cause)
}
func ServiceUnavailable(msg string, cause error) *AppError {
	return New(ErrUnavailable, http.StatusServiceUnavailable, msg, cause)
}
