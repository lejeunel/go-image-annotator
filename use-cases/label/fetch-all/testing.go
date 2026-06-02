package fetchall

import (
	"context"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessFetchLabels(r Response) {
	p.GotSuccess = true
	p.Got = r
}

type FakeRepo struct {
	Labels     []string
	Count_     int64
	ErrOnCount bool
	ErrOnFetch bool
	Err        error
}

func (r *FakeRepo) Count() (int64, error) {
	if r.ErrOnCount {
		return 0, r.Err
	}
	return r.Count_, nil
}

func (r *FakeRepo) FetchAll() ([]string, error) {
	if r.ErrOnFetch {
		return nil, r.Err
	}
	return r.Labels, nil
}

type FailingAuth struct {
}

func (f FailingAuth) FetchAllLabels(ctx context.Context) error {
	return e.ErrAuth
}
