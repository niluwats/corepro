package errs

import "net/http"

var (
	Wrong_pw        = "incorrect password"
	Wrong_em        = "incorrect email"
	Exist_user      = "user with this email address already exists"
	Unverified_em   = "please verify your email to login"
	Deactivated_acc = "your account is deactivated"
	Required        = "enter all required field"
	Invalid_pw      = "Your password should be more than 6 characters long with including atleast one special character,capital letter and a number "
)

type AppError struct {
	Code    int    `json:",omitempty"`
	Message string `json:"message"`
}

//	AsMessage returns error message as a pointer of AppError
func (e AppError) AsMessage() *AppError {
	return &AppError{
		Message: e.Message,
	}
}

//	NewUnexpectedError returns 404 http.StatusNotFound error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

//	NewUnexpectedError returns 500 http.StatusInternalServerError error
func NewUnexpectedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

//	NewValidationError returns 422 http.StatusUnprocessableEntity error
func NewValidationError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusUnprocessableEntity,
	}
}

//	NewAuthenticationError returns 401 http.StatusUnauthorized error
func NewAuthenticationError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

//	NewAuthorizationError returns 403 http.StatusForbidden error
func NewAuthorizationError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusForbidden,
	}
}
