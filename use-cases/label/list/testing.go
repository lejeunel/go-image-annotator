package list

import (
	"context"
	l "github.com/lejeunel/go-image-annotator/entities/label"
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
	ErrOnCount bool
	ErrOnList  bool
	Err        error
	Count_     int
}

func (r *FakeRepo) Count() (int64, error) {
	if r.ErrOnCount {
		return 0, r.Err
	}
	return int64(r.Count_), nil
}

func (r *FakeRepo) List(req Request) ([]*l.Label, error) {
	if r.ErrOnList {
		return nil, r.Err

	}

	result := []*l.Label{}
	for range req.PageSize {
		l := l.NewLabel(l.NewLabelId(), "a-label")
		result = append(result, &l)
	}
	return result, nil

}

type FailingAuth struct {
}

func (f FailingAuth) ListLabels(ctx context.Context) error {
	return e.ErrAuth
}
