package validation

import (
	"fmt"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"regexp"
)

type Validator interface {
	Validate(string) error
}

type FakeNameValidator struct {
	Err error
}

func (v *FakeNameValidator) Validate(name string) error {
	return v.Err
}

type NameValidator struct {
}

var validName = regexp.MustCompile(`^[a-z0-9-]+$`)

func (v *NameValidator) Validate(name string) error {
	if !validName.MatchString(name) {
		return fmt.Errorf("checking for illegal characters (capital letters, special characters except '-'): %w", e.ErrValidation)
	}
	return nil

}

func NewNameValidator() *NameValidator {
	return &NameValidator{}
}
