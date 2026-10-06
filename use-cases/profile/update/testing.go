package update

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        pr.Profile
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessUpdateProfile(r pr.Profile) {
	p.GotSuccess = true
	p.Got = r
}
