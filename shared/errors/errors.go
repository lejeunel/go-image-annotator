package errors

import (
	"errors"
)

var ErrInternal = errors.New("internal error")
var ErrDuplicate = errors.New("duplicate resource error")
var ErrNotFound = errors.New("resource not found error")
var ErrDependency = errors.New("dependency error")
var ErrValidation = errors.New("validation error")
var ErrImageFormat = errors.New("forbidden image format")
var ErrURLParsing = errors.New("url parsing error")
var ErrLabelLimitExceeded = errors.New("label limit count exceeded error")
var ErrAuthentication = errors.New("authentication error")
var ErrAuthorization = errors.New("authorization error")
var ErrPrincipalProvider = errors.New("error extracting principal identity")
var ErrPasswordMismatch = errors.New("password mismatch error")
var ErrInvalidPassword = errors.New("invalid password error")
var ErrExpiredToken = errors.New("expired token error")
var ErrForbiddenOp = errors.New("forbidden operation error")
