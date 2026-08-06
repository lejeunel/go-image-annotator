package list

import (
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        []m.MetaData
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessListMetadata(r []m.MetaData) {
	p.Got = r
	p.GotSuccess = true
}
