package update

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	group := "a-group"
	itr := New(&FakeCollectionRepo{},
		&FakeGroupRepo{Return: &group},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestUpdateNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := New(&FakeCollectionRepo{}, &FakeGroupRepo{})
	itr.Execute(t.Context(), Request{Name: non_existing_name, NewName: "new-name"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateCollection(t *testing.T) {
	name := "name"
	p := &FakePresenter{}
	repo := &FakeCollectionRepo{Names: []string{name}}
	itr := New(repo, &FakeGroupRepo{})
	req := Request{Name: name,
		NewName:        "updated-name",
		NewDescription: "updated-description"}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.NewName, p.Got.Name)
	assert.Equal(t, req.NewDescription, p.Got.Description)
}

func TestUpdateCollectionWithNameAlreadyTakenShouldFail(t *testing.T) {

	p := &FakePresenter{}
	name := "name"
	existing_name := "existing-name"
	itr := New(&FakeCollectionRepo{Names: []string{name, existing_name}},
		&FakeGroupRepo{})
	itr.Execute(t.Context(), Request{Name: name, NewName: existing_name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateCollectionWithUnchangedNameShouldSucceed(t *testing.T) {

	p := &FakePresenter{}
	name := "name"
	itr := New(&FakeCollectionRepo{Names: []string{name}}, &FakeGroupRepo{})
	itr.Execute(t.Context(), Request{Name: name, NewName: name}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	name := "name"
	itr := New(&FakeCollectionRepo{Names: []string{name},
		Err: e.ErrInternal}, &FakeGroupRepo{})
	itr.Execute(t.Context(),
		Request{Name: name, NewName: name}, p)
	assert.True(t, p.GotInternalErr)
}
