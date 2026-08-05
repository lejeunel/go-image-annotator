package fake

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type ValueValidator struct {
	Invalid bool
}

func (v *ValueValidator) Validate(value any) error {
	if v.Invalid {
		return e.ErrValidation
	}
	return nil
}

type MetaDataRepo struct {
	AddedKey   string
	AddedValue any
	ErrOnAdd   error
}

func (r *MetaDataRepo) Add(key string, value any) error {
	if r.ErrOnAdd != nil {
		return r.ErrOnAdd
	}
	r.AddedKey = key
	r.AddedValue = value
	return nil

}
