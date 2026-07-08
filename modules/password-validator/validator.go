package password_validator

import (
	gpv "github.com/wagslane/go-password-validator"
)

type PasswordValidator interface {
	Validate(string) error
}

type MyPasswordValidator struct {
	minEntropyBits int
}

func (pv MyPasswordValidator) Validate(password string) error {
	return gpv.Validate(password, float64(pv.minEntropyBits))
}

func New(minEntropyBits int) MyPasswordValidator {
	return MyPasswordValidator{minEntropyBits}
}
