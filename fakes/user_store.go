package fake

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type UserStore struct {
	ErrOnFind     error
	ErrOnList     error
	Return        *u.User
	ExistingUsers []u.User
}

func (s UserStore) Find(id u.UserId) (*u.User, error) {
	if s.ErrOnFind != nil {
		return nil, s.ErrOnFind
	}
	return s.Return, nil
}

func (s UserStore) List(p pag.PaginationParams) ([]u.User, *pag.Pagination, error) {
	if s.ErrOnList != nil {
		return nil, nil, s.ErrOnList
	}
	return s.ExistingUsers, &pag.Pagination{
		PageSize:     p.PageSize,
		Page:         p.Page,
		TotalRecords: int64(len(s.ExistingUsers)),
		TotalPages:   int64(len(s.ExistingUsers) / p.PageSize),
	}, nil
}
