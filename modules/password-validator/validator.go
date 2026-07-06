package password_validator

import (
	gpv "github.com/wagslane/go-password-validator"
)

type PasswordValidator struct {
	minEntropyBits int
}

func (pv PasswordValidator) Validate(password string) error {
	return gpv.Validate(password, float64(pv.minEntropyBits))
}

func New(minEntropyBits int) PasswordValidator {
	return PasswordValidator{minEntropyBits}
}
