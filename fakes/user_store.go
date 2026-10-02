package fake

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type UserStore struct {
	ErrOnFind error
	Return    *u.User
}

func (s UserStore) Find(id u.UserId) (*u.User, error) {
	if s.ErrOnFind != nil {
		return nil, s.ErrOnFind
	}
	return s.Return, nil
}
