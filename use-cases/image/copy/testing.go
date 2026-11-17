package copy

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessCopyImage(r Response) {
	p.Got = r
	p.GotSuccess = true
}
