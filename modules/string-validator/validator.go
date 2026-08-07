package validator

import (
	"fmt"
	"regexp"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Validator interface {
	Validate(string) error
}

type NameValidator struct{}

var validName = regexp.MustCompile(`^[a-z0-9-]+$`)

func (v *NameValidator) Validate(name string) error {
	if !validName.MatchString(name) {
		return fmt.Errorf(
			"checking for illegal characters (capital letters, special characters except '-'): %w",
			e.ErrValidation,
		)
	}
	return nil
}

func NewNameValidator() *NameValidator {
	return &NameValidator{}
}
