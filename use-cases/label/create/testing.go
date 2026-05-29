package create

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
	"slices"
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
	Err   error
	Names []string
	Got   lbl.Label
}

func (r *FakeRepo) Create(l lbl.Label) error {
	if r.Err != nil {
		return r.Err
	}
	r.Got = l
	return nil
}
func (r *FakeRepo) Exists(name string) (bool, error) {
	if slices.Contains(r.Names, name) {
		return true, nil
	}
	return false, nil
}
