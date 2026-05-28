package create

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	st "github.com/lejeunel/go-image-annotator-v2/shared/testing"
	v "github.com/lejeunel/go-image-annotator-v2/shared/validation"
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
	st.AssertEqual(t, "name", repo.Got.Name, req.Name)
	st.AssertEqual(t, "description", repo.Got.Description, req.Description)
	st.AssertEqual(t, "creation date", repo.Got.CreatedAt, now)
	st.AssertEqual(t, "id", repo.Got.Id.IsNil(), false)
}
