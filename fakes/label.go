package fake

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
	"slices"
)

type LabelRepo struct {
	ErrOnFind     error
	ErrOnCreate   error
	ErrOnDelete   error
	ErrOnIsUsed   error
	ErrOnExists   error
	ErrOnFetch    error
	ErrOnList     error
	ErrOnUpdate   error
	ErrOnCount    error
	Return        lbl.Label
	FetchedName   string
	ExistingNames []string
	Created       lbl.Label
	IsUsed_       bool
	Count_        int64
	GotUpdatable  lbl.UpdatableModel
}

func (r *LabelRepo) FindLabel(name string) (*lbl.Label, error) {
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}
	r.FetchedName = name
	return &r.Return, nil
}

func (r *LabelRepo) Create(l lbl.Label) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}
	r.Created = l
	return nil
}
func (r *LabelRepo) Exists(name string) (bool, error) {
	if r.ErrOnExists != nil {
		return false, r.ErrOnExists
	}
	if slices.Contains(r.ExistingNames, name) {
		return true, nil
	}
	return false, nil
}

func (r *LabelRepo) Delete(string) error {
	if r.ErrOnDelete != nil {
		return r.ErrOnDelete
	}
	return nil
}

func (r *LabelRepo) IsUsed(n string) (*bool, error) {
	res := true
	if r.ErrOnIsUsed != nil {
		return nil, r.ErrOnIsUsed
	}
	if r.IsUsed_ {
		return &res, nil

	}
	res = false
	return &res, nil
}

func (r *LabelRepo) Count() (int64, error) {
	if r.ErrOnCount != nil {
		return 0, r.ErrOnCount
	}
	return r.Count_, nil
}

func (r *LabelRepo) FetchAll() ([]string, error) {
	if r.ErrOnFetch != nil {
		return nil, r.ErrOnFetch
	}
	return r.ExistingNames, nil
}

func (r *LabelRepo) List(req pag.PaginationParams) ([]*lbl.Label, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	result := []*lbl.Label{}
	for range req.PageSize {
		l := lbl.NewLabel(lbl.NewLabelId(), "a-label")
		result = append(result, &l)
	}
	return result, nil

}

func (r *LabelRepo) Update(m lbl.UpdatableModel) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdatable = m
	return nil
}
