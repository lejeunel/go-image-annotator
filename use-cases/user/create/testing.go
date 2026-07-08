package create

import (
	"context"
	"slices"

	tk "github.com/lejeunel/go-image-annotator/entities/token"
	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(r Response) {
	p.GotSuccess = true
	p.Got = r
}

type FakeRepo struct {
	Err error
	Ids []string
	Got *usr.User
}

func (r *FakeRepo) Create(u usr.User) error {
	if r.Err != nil {
		return r.Err
	}
	r.Got = &u
	return nil
}
func (r *FakeRepo) Exists(id string) (bool, error) {
	if slices.Contains(r.Ids, id) {
		return true, nil
	}
	return false, nil
}

type FailingAuth struct {
}

func (f FailingAuth) CreateUser(ctx context.Context) error {
	return e.ErrAuthorization
}

type FakeTokenGenerator struct {
	Value                  string
	Hash_                  []byte
	GeneratedHashFromValue string
}

func (t FakeTokenGenerator) Generate() (*tk.Token, error) {
	return &tk.Token{Value: t.Value, Hash: t.Hash_}, nil
}

func (t *FakeTokenGenerator) Hash(value string) []byte {
	t.GeneratedHashFromValue = value
	return t.Hash_
}
