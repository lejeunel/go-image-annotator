package fake

import (
	"slices"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
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
	ExistingKeys      []string
	AddedKey          string
	DeletedKey        string
	AddedValue        any
	UpdatedKey        string
	UpdatedValue      any
	ErrOnAdd          error
	ErrOnKeyExists    error
	ErrOnDelete       error
	ErrOnList         error
	ErrOnGet          error
	ErrOnUpdate       error
	ErrOnDeleteAll    error
	ReturnList        []m.MetaData
	ReturnValue       any
	DeletedAllImageId *im.ImageId
}

func (r *MetaDataRepo) KeyExists(n clc.CollectionName, id im.ImageId, key string) (bool, error) {
	if r.ErrOnKeyExists != nil {
		return false, r.ErrOnKeyExists
	}
	exist := true
	if slices.Contains(r.ExistingKeys, key) {
		return exist, nil
	}
	exist = false
	return exist, nil
}

func (r *MetaDataRepo) Add(n clc.CollectionName, id im.ImageId, key string, value any) error {
	if r.ErrOnAdd != nil {
		return r.ErrOnAdd
	}
	r.AddedKey = key
	r.AddedValue = value
	return nil
}

func (r *MetaDataRepo) Delete(n clc.CollectionName, id im.ImageId, key string) error {
	if r.ErrOnDelete != nil {
		return r.ErrOnDelete
	}
	r.DeletedKey = key

	return nil
}
func (r *MetaDataRepo) DeleteAll(n clc.CollectionName, id im.ImageId) error {
	if r.ErrOnDeleteAll != nil {
		return r.ErrOnDeleteAll
	}
	r.DeletedAllImageId = &id

	return nil
}

func (r *MetaDataRepo) List(clc.CollectionName, im.ImageId) ([]m.MetaData, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}
	return r.ReturnList, nil
}

func (r *MetaDataRepo) GetValue(clc.CollectionName, im.ImageId, string) (*any, error) {
	if r.ErrOnGet != nil {
		return nil, r.ErrOnGet
	}
	return &r.ReturnValue, nil
}

func (r *MetaDataRepo) UpdateValue(
	c clc.CollectionName,
	i im.ImageId,
	key string,
	value any,
) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.UpdatedKey = key
	r.UpdatedValue = value
	return nil
}
