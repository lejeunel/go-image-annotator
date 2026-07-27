package find

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
	"io"
)

type FakePresenter struct {
	Got        []byte
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessReadPolicy(r io.Reader) {
	p.GotSuccess = true
	p.Got, _ = io.ReadAll(r)
}
