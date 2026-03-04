package errors

type AppError struct {
	Code	ErrorCode
	Message string
	HTTPStatus int
	Err 	 error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NotFound(code ErrorCode, message string) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		HTTPStatus: 404,
	}
}

func Unauthorized(code ErrorCode, message string) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		HTTPStatus: 401,
	}
}

func Conflict(code ErrorCode, message string) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		HTTPStatus: 409,
	}
}

func Validation(code ErrorCode, message string) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		HTTPStatus: 400,
	}
}

func Internal(message string, err error) *AppError {
	return &AppError{
		Code: ErrInternal,
		Message: message,
		Err: err,
		HTTPStatus: 500,
	}
}