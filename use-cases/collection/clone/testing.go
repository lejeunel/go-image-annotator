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

type TestingTransactor struct {
	Repos
}

func (m *TestingTransactor) RunInTx(
	fn func(Repos) error) error {
	return fn(m.Repos)
}

func NewTestingCloner() Interactor {
	repos := Repos{&fk.ImageRepo{}, &fk.CollectionRepo{}, &fk.AnnotationRepo{}}
	return New(
		repos,
		&TestingTransactor{repos},
		&fk.GroupRepo{},
		&fk.ImageStore{}, &fk.EventLogger{}, fk.NewLogger(), &fk.JobQueue{})
}
