package errors

import (
	"errors"
)

var (
	ErrInternal           = errors.New("internal error")
	ErrDuplicate          = errors.New("duplicate resource error")
	ErrNotFound           = errors.New("resource not found error")
	ErrDependency         = errors.New("dependency error")
	ErrValidation         = errors.New("validation error")
	ErrImageFormat        = errors.New("forbidden image format")
	ErrURLParsing         = errors.New("url parsing error")
	ErrLabelLimitExceeded = errors.New("label limit count exceeded error")
	ErrAuthentication     = errors.New("authentication error")
	ErrAuthorization      = errors.New("authorization error")
	ErrPrincipalProvider  = errors.New("error extracting principal identity")
	ErrPasswordMismatch   = errors.New("password mismatch error")
	ErrInvalidPassword    = errors.New("invalid password error")
	ErrExpiredToken       = errors.New("expired token error")
	ErrForbiddenOp        = errors.New("forbidden operation error")
)
