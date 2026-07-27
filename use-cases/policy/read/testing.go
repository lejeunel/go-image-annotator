package read

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        string
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessReadPolicy(policies string) {
	p.GotSuccess = true
	p.Got = policies
}
