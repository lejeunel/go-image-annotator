package fake

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type StringValidator struct {
	Invalid   bool
	Validated string
}

func (v *StringValidator) Validate(s string) error {
	v.Validated = s
	if v.Invalid {
		return e.ErrValidation
	}

	return nil
}
