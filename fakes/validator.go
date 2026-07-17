package fake

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Validator struct {
	Invalid bool
}

func (v *Validator) Validate(string) error {
	if v.Invalid {
		return e.ErrInvalidPassword
	}
	return nil
}
