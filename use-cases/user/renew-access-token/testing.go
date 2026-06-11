package renew_token

import (
	"context"

	tok "github.com/lejeunel/go-image-annotator/app/token"
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
	Err     error
	GotId   usr.UserId
	GotHash []byte
	Missing bool
}

func (r *FakeRepo) SetAccessTokenHash(id usr.UserId, hash []byte) error {
	if r.Err != nil {
		return r.Err
	}
	r.GotId = id
	r.GotHash = hash
	return nil
}
func (r *FakeRepo) Exists(id string) (bool, error) {
	if r.Missing {
		return false, nil
	}
	return true, nil
}

type FailingAuth struct {
}

func (f FailingAuth) RenewToken(ctx context.Context, id string) error {
	return e.ErrAuth
}

type FakeTokenGenerator struct {
	Token string
	Hash_ []byte
}

func (t *FakeTokenGenerator) Generate() (*tok.TokenPair, error) {
	return &tok.TokenPair{Token: t.Token, Hash: t.Hash_}, nil
}
