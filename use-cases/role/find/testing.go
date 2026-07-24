package find

import (
	r "github.com/lejeunel/go-image-annotator/entities/role"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        r.Role
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessFindRole(role r.Role) {
	p.GotSuccess = true
	p.Got = role
}
