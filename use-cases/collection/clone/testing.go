package clone

import (
	"context"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	testing "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeCloner struct {
	Err     error
	GotTask t.CloneTask
}

func (c *FakeCloner) Clone(task t.CloneTask) error {
	if c.Err != nil {
		return c.Err
	}
	c.GotTask = task
	return nil
}

type FakePresenter struct {
	GotTask    *t.TaskSpecs
	GotSuccess bool
	testing.TestingErrPresenter
}

func (p *FakePresenter) SuccessSubmitCloneTask(t t.TaskSpecs) {
	p.GotTask = &t
	p.GotSuccess = true
}

type FailingAuth struct {
}

func (f FailingAuth) CloneCollection(ctx context.Context, g string) error {
	return e.ErrAuthorization
}
