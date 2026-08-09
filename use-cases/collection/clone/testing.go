package clone

import (
	fk "github.com/lejeunel/go-image-annotator/fakes"
	testing "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	testing.TestingErrPresenter
}

func (p *FakePresenter) SuccessSubmitCloneTask(r Response) {
	p.Got = r
	p.GotSuccess = true
}

func NewTestingCloner() Interactor {
	return New(
		&fk.ImageStore{},
		&fk.ImageRepo{},
		&fk.CollectionRepo{},
		&fk.GroupRepo{},
		&fk.EventLogger{}, fk.NewLogger(), &fk.JobQueue{})
}
