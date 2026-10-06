package create

import (
	c "github.com/lejeunel/go-image-annotator/entities/collection"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        c.Collection
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessCreateCollection(r c.Collection) {
	p.GotSuccess = true
	p.Got = r
}
