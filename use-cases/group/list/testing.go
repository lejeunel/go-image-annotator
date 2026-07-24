package list

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        []grp.Group
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessListGroups(groups []grp.Group) {
	p.GotSuccess = true
	p.Got = groups
}
