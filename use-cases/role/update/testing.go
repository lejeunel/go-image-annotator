package update

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessUpdateRole(r Response) {
	p.GotSuccess = true
	p.Got = r
}

type FailingAuth struct {
}
