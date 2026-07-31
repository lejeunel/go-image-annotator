package find

import (
	ta "github.com/lejeunel/go-image-annotator/entities/task"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        ta.Task
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessFindTask(t ta.Task) {
	p.GotSuccess = true
	p.Got = t
}
