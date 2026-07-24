package list

import (
	r "github.com/lejeunel/go-image-annotator/entities/role"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        []r.Role
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessListRoles(roles []r.Role) {
	p.GotSuccess = true
	p.Got = roles
}
