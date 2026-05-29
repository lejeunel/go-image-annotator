package update

import (
	"slices"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Names []string
	Got   Model
}

func (r *FakeRepo) Update(m Model) error {
	r.Got = m
	return nil
}
func (r *FakeRepo) Exists(n string) (bool, error) {
	if slices.Contains(r.Names, n) {
		return true, nil
	}
	return false, nil
}

type FakeErrRepo struct {
	err error
}

func (r *FakeErrRepo) Update(m Model) error {
	return e.ErrInternal
}

func (r *FakeErrRepo) Exists(n string) (bool, error) {
	return false, r.err
}

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(r Response) {
	p.GotSuccess = true
	p.Got = r
}
