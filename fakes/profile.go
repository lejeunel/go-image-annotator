package fake

import (
	"slices"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type ProfileRepo struct {
	ErrOnCreate    error
	ErrOnExists    error
	ErrOnAddLabel  error
	ErrOnFind      error
	ErrOnDelete    error
	ErrOnIsUsed    error
	ErrOnList      error
	ErrOnCount     error
	ErrOnGetGroup  error
	ErrOnUpdate    error
	ExistingNames  []string
	Created        []pr.Profile
	AddedLabels    []lbl.LabelName
	Return         pr.Profile
	ReturnGroup    string
	IsUsed_        bool
	Count_         int
	GotUpdateModel pr.UpdateModel
}

func (r *ProfileRepo) Create(p pr.Profile) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}

	r.Created = append(r.Created, p)
	return nil
}

func (r *ProfileRepo) Exists(name string) (*bool, error) {
	if r.ErrOnExists != nil {
		return nil, r.ErrOnExists
	}

	exist := true
	if slices.Contains(r.ExistingNames, name) {
		return &exist, nil
	}
	exist = false
	return &exist, nil
}

func (r *ProfileRepo) AddLabel(profileName pr.ProfileName, label lbl.LabelName) error {
	if r.ErrOnAddLabel != nil {
		return r.ErrOnAddLabel
	}
	r.AddedLabels = append(r.AddedLabels, label)
	return nil
}

func (r *ProfileRepo) Find(name string) (*pr.Profile, error) {
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}

	return &r.Return, nil
}

func (r *ProfileRepo) Delete(name string) error {
	if r.ErrOnDelete != nil {
		return r.ErrOnDelete
	}

	return nil
}

func (r *ProfileRepo) IsUsed(n string) (*bool, error) {
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

func (r *ProfileRepo) Count() (*int64, error) {
	count := int64(0)
	if r.ErrOnCount != nil {
		return &count, r.ErrOnCount
	}
	res := int64(r.Count_)
	return &res, nil
}

func (r *ProfileRepo) List(req pa.PaginationParams) ([]pr.Profile, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	result := []pr.Profile{}
	for range req.PageSize {
		result = append(result, r.Return)
	}
	return result, nil
}

func (r *ProfileRepo) ListAll() ([]string, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	result := []string{r.Return.Name}
	return result, nil
}

func (r *ProfileRepo) GetGroup(name string) (*string, error) {
	if r.ErrOnGetGroup != nil {
		return nil, r.ErrOnGetGroup
	}
	return &r.ReturnGroup, nil
}

func (r *ProfileRepo) Update(m pr.UpdateModel) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdateModel = m
	return nil
}
