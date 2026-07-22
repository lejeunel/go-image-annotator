package find

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        u.User
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessFindUser(user u.User) {
	p.GotSuccess = true
	p.Got = user
}
