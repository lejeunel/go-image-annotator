package create

import (
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	st "github.com/lejeunel/go-image-annotator-v2/shared/testing"
	v "github.com/lejeunel/go-image-annotator-v2/shared/validation"
	"testing"
)

func TestCreateLabelWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-label"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Names: []string{name}}, &v.FakeNameValidator{})
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
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal}, &v.FakeNameValidator{})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr {
		t.Fatal("expected internal error, but got none")
	}
}

func TestCreateLabelWithInvalidNameShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &v.FakeNameValidator{Err: e.ErrValidation})
	itr.Execute(Request{Name: "invalid-name"}, p)
	if !p.GotValidationErr {
		t.Fatal("expected validation error, but go none")
	}
}

func TestCreateLabel(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := NewInteractor(repo, &v.FakeNameValidator{})
	name := "a-name"
	desc := "a-description"
	req := Request{Name: name, Description: desc}
	itr.Execute(req, p)

	st.AssertEqual(t, "name", p.Got.Name, name)
	st.AssertEqual(t, "description", p.Got.Description, desc)
	st.AssertEqual(t, "name", repo.Got.Name, name)
	st.AssertEqual(t, "description", repo.Got.Description, desc)
	st.AssertEqual(t, "id", repo.Got.Id.IsNil(), false)
}
