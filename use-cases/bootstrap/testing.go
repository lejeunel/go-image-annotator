package bootstrap

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	Got        Response
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessBootstrap(r Response) {
	p.GotSuccess = true
	p.Got = r
}
