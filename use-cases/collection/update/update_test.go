package update

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

func TestUpdateNonExistingCollectionShouldFail(t *testing.T) {

	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(Request{Name: non_existing_name, NewName: "new-name"}, p)
	if !p.GotNotFoundErr {
		t.Fatal("expected not found error, but got none")
	}
	if p.GotSuccess {
		t.Fatal("expected no success")
	}
}

func TestUpdateCollection(t *testing.T) {
	name := "name"
	new_name := "updated-name"
	new_description := "updated-description"

	p := &FakePresenter{}
	repo := &FakeRepo{Names: []string{name}}
	itr := NewInteractor(repo)
	req := Request{Name: name, NewName: new_name, NewDescription: new_description}
	wantr := Model{Name: name, NewName: new_name, NewDescription: new_description}
	wantp := Response{Name: new_name, Description: new_description}
	itr.Execute(req, p)
	if p.Got != wantp {
		t.Fatalf("expected %v, got %v", wantp, p.Got)

	}
	if repo.Got != wantr {
		t.Fatalf("expected %v, got %v", req, repo.Got)

	}
}

func TestUpdateCollectionWithNameAlreadyTakenShouldFail(t *testing.T) {

	p := &FakePresenter{}
	name := "name"
	existing_name := "existing-name"
	itr := NewInteractor(&FakeRepo{Names: []string{name, existing_name}})
	itr.Execute(Request{Name: name, NewName: existing_name}, p)
	if !p.GotDuplicationErr {
		t.Fatal("expected duplication error, but got none")
	}
	if p.GotSuccess {
		t.Fatal("expected no success")
	}
}

func TestUpdateCollectionWithUnchangedNameShouldSucceed(t *testing.T) {

	p := &FakePresenter{}
	name := "name"
	itr := NewInteractor(&FakeRepo{Names: []string{name}})
	itr.Execute(Request{Name: name, NewName: name}, p)
	if !p.GotSuccess {
		t.Fatal("expected success")
	}
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeErrRepo{e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr {
		t.Fatal("expected internal error, but got none")
	}
}
