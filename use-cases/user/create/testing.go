package create

import (
	"context"
	"slices"

	g "github.com/lejeunel/go-image-annotator/app/token-generator"
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
	return e.ErrAuth
}

type FakeTokenGenerator struct {
	Token string
	Hash_ []byte
}

func (t *FakeTokenGenerator) Generate() (*g.TokenPair, error) {
	return &g.TokenPair{Token: t.Token, Hash: t.Hash_}, nil
}
