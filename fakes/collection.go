package fake

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"slices"
)

type CollectionRepo struct {
	ErrOnCreate    error
	ErrOnExists    error
	ErrOnFind      error
	ErrOnDelete    error
	ErrOnCount     error
	ErrOnList      error
	ErrOnUpdate    error
	ExistingNames  []string
	IsPopulated_   bool
	Return         clc.Collection
	Count_         int
	Got            clc.Collection
	GotUpdateModel clc.UpdateModel
}

func (r *CollectionRepo) Create(c clc.Collection) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}

	r.Got = c
	return nil
}

func (r *CollectionRepo) Exists(name string) (bool, error) {
	if r.ErrOnExists != nil {
		return false, r.ErrOnExists
	}
	if slices.Contains(r.ExistingNames, name) {
		return true, nil
	}
	return false, nil
}

func (r *CollectionRepo) FindCollectionByName(name string) (*clc.Collection, error) {
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}

	return &r.Return, nil

}
func (r *CollectionRepo) Delete(string) error {

	if r.ErrOnDelete != nil {
		return r.ErrOnDelete
	}
	return nil
}

func (r *CollectionRepo) IsPopulated(c string) (*bool, error) {
	res := true
	if r.IsPopulated_ {
		return &res, nil
	}
	res = false
	return &res, nil
}

func (r *CollectionRepo) Count() (*int64, error) {
	count := int64(0)
	if r.ErrOnCount != nil {
		return &count, r.ErrOnCount
	}
	res := int64(r.Count_)
	return &res, nil
}

func (r *CollectionRepo) List(req pa.PaginationParams) ([]*clc.Collection, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	result := []*clc.Collection{}
	for range req.PageSize {
		result = append(result, &r.Return)
	}
	return result, nil
}

func (r *CollectionRepo) Update(m clc.UpdateModel) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdateModel = m
	return nil
}
