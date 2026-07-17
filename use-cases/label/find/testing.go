package find

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        l.Label
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessFindLabel(l l.Label) {
	p.GotSuccess = true
	p.Got = l
}
