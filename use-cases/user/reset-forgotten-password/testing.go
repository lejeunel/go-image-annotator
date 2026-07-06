package reset_forgotten_password

import (
	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success() {
	p.GotSuccess = true
}

type FakeRepo struct {
	GotId       usr.UserId
	GotHash     []byte
	Missing     bool
	Return      *usr.ForgotPasswordState
	ErrOnUpdate error
}

func (r *FakeRepo) FindForgottenPassword(hash []byte) (*usr.ForgotPasswordState, error) {
	if r.Missing {
		return nil, e.ErrNotFound
	}
	r.GotHash = hash
	return r.Return, nil
}
func (r *FakeRepo) UpdatePassword(id usr.UserId, hash []byte) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotId = id
	r.GotHash = hash
	return nil
}

type FakeTokenHasher struct {
	GotToken   string
	ReturnHash []byte
}

func (t *FakeTokenHasher) Hash(token string) []byte {
	t.GotToken = token
	return t.ReturnHash
}

type FakePasswordValidator struct {
	Invalid bool
}

func (v *FakePasswordValidator) Validate(string) error {
	if v.Invalid {
		return e.ErrInvalidPassword
	}
	return nil
}
