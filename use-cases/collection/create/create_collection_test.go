package create

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
	"github.com/stretchr/testify/assert"
)

func TestCreateCollectionWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Names: []string{name}})
	itr.Execute(Request{Name: name}, p)
	if !p.GotDuplicationErr {
		t.Fatal("expected duplication error, but go none")
	}
	if p.GotSuccess {
		t.Fatal("expected no success")
	}
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal},
		WithNameValidator(&v.FakeNameValidator{}))
	itr.Execute(Request{}, p)
	if !p.GotInternalErr {
		t.Fatal("expected internal error, but got none")
	}
}

func TestCreateCollectionWithInvalidNameShouldFail(t *testing.T) {
	name := "my-collection%/"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Names: []string{name}},
		WithNameValidator(&v.FakeNameValidator{Err: e.ErrValidation}))
	itr.Execute(Request{Name: name}, p)
	if !p.GotValidationErr {
		t.Fatal("expected validation error, but go none")
	}
}

func TestCreateCollection(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	now := time.Now()
	itr := NewInteractor(repo, WithClock(clockwork.NewFakeClockAt(now)))
	req := Request{Name: "a-name", Description: "a-descriptin"}
	itr.Execute(req, p)
	assert.Equal(t, repo.Got.Name, req.Name, "name")
	assert.Equal(t, repo.Got.Description, req.Description, "description")
	assert.Equal(t, repo.Got.CreatedAt, now, "creation date")
	assert.Equal(t, repo.Got.Id.IsNil(), false, "id")
}
