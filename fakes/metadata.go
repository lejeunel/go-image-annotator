package fake

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"slices"
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
	ExistingKeys   []string
	AddedKey       string
	AddedValue     any
	ErrOnAdd       error
	ErrOnKeyExists error
}

func (r *MetaDataRepo) KeyExists(key string) (*bool, error) {
	if r.ErrOnKeyExists != nil {
		return nil, r.ErrOnKeyExists
	}
	var exist = true
	if slices.Contains(r.ExistingKeys, key) {
		return &exist, nil
	}
	exist = false
	return &exist, nil

}

func (r *MetaDataRepo) Add(key string, value any) error {
	if r.ErrOnAdd != nil {
		return r.ErrOnAdd
	}
	r.AddedKey = key
	r.AddedValue = value
	return nil

}
