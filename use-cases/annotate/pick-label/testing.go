package pick

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        []string
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessFetchLabels(labels []string) {
	p.GotSuccess = true
	p.Got = labels
}
