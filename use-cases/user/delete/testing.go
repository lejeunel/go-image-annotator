package delete

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessDeleteUser(u.UserId) {
	p.GotSuccess = true
}
