package list

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        []string
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessListAllProfiles(r []string) {
	p.GotSuccess = true
	p.Got = r
}
