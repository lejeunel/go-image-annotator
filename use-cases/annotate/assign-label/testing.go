package assign_label

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessAddLabel(r Response) {
	p.Got = r
	p.GotSuccess = true
}
