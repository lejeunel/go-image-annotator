package create

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        rl.Role
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessCreateRole(r rl.Role) {
	p.GotSuccess = true
	p.Got = r
}
