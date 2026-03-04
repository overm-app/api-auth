package errors

type ErrorCode string

const (
	// Auth
    ErrInvalidCredentials  ErrorCode = "INVALID_CREDENTIALS"
    ErrUnauthorized        ErrorCode = "UNAUTHORIZED"
    ErrTokenExpired        ErrorCode = "TOKEN_EXPIRED"
    ErrTokenInvalid        ErrorCode = "TOKEN_INVALID"
    ErrRefreshTokenInvalid ErrorCode = "INVALID_REFRESH_TOKEN"
    ErrAuthProviderConflict ErrorCode = "AUTH_PROVIDER_CONFLICT"
    ErrOAuthCodeInvalid    ErrorCode = "OAUTH_CODE_INVALID"

    // Resource
    ErrNotFound            ErrorCode = "NOT_FOUND"
    ErrAlreadyExists       ErrorCode = "ALREADY_EXISTS"
    ErrEmailAlreadyExists  ErrorCode = "EMAIL_ALREADY_EXISTS"

	// Conflict
	ErrConflict            ErrorCode = "CONFLICT"

    // Validation
    ErrValidation          ErrorCode = "VALIDATION_ERROR"

    // Internal
    ErrInternal            ErrorCode = "INTERNAL_ERROR"
)