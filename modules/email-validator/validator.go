package email_validator

import (
	"fmt"
	"net/mail"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type EmailValidator interface {
	Validate(string) error
}

type MyEmailValidator struct{}

func (ev MyEmailValidator) Validate(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("validation email address %v: %w: %w", email, err, e.ErrValidation)
	}
	return nil
}

func New() MyEmailValidator {
	return MyEmailValidator{}
}
