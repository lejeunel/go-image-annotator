package create

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        pr.Profile
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessCreateProfile(profile pr.Profile) {
	p.GotSuccess = true
	p.Got = profile
}
