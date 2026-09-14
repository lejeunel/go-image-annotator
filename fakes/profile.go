package fake

import (
	"slices"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type CreatedProfile struct {
	Id          pr.ProfileId
	Name        pr.ProfileName
	Description pr.ProfileDescription
}

type ProfileRepo struct {
	ErrOnCreate   error
	ErrOnExists   error
	ErrOnAddLabel error
	ErrOnFind     error
	ErrOnDelete   error
	ErrOnIsUsed   error
	ErrOnList     error
	ErrOnCount    error
	ExistingNames []string
	Created       []CreatedProfile
	AddedLabels   []lbl.LabelName
	Return        pr.Profile
	IsUsed_       bool
	Count_        int
}

func (r *ProfileRepo) Create(
	id pr.ProfileId,
	name pr.ProfileName,
	desc pr.ProfileDescription,
) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}

	r.Created = append(r.Created,
		CreatedProfile{Id: id, Name: name, Description: desc})
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
