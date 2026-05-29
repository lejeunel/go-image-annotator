package list

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
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
		result = append(result, l.NewLabel(l.NewLabelId(), "a-label"))
	}
	return result, nil

}
