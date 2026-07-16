package find

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        clc.Collection
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessFindCollection(c clc.Collection) {
	p.GotSuccess = true
	p.Got = c
}
