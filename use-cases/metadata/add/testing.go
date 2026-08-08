package add

import (
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessAddMetadata(m.MetaData) {
	p.GotSuccess = true
}
