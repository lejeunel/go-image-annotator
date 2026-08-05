package fake

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type StringValidator struct {
	Invalid bool
}

func (v *StringValidator) Validate(string) error {
	if v.Invalid {
		return e.ErrValidation
	}
	return nil
}
