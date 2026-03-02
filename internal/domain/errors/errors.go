package errors

type AppError struct {
	Code	string
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

func NotFound(code, message string) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		HTTPStatus: 404,
	}
}

func Unauthorized(code, message string) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		HTTPStatus: 401,
	}
}

func Conflict(code, message string) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		HTTPStatus: 409,
	}
}

func Validation(code, message string) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		HTTPStatus: 400,
	}
}

func Internal(message string, err error) *AppError {
	return &AppError{
		Code: "INTERNAL_ERROR",
		Message: message,
		Err: err,
		HTTPStatus: 500,
	}
}