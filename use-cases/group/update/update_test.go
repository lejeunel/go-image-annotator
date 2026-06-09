package update

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestUpdateNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := New(&FakeRepo{})
	itr.Execute(t.Context(), Request{Name: non_existing_name, NewName: "new-name"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateCollection(t *testing.T) {
	name := "name"
	p := &FakePresenter{}
	repo := &FakeRepo{Names: []string{name}}
	itr := New(repo)
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
	itr := New(&FakeRepo{Names: []string{name, existing_name}})
	itr.Execute(t.Context(), Request{Name: name, NewName: existing_name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateCollectionWithUnchangedNameShouldSucceed(t *testing.T) {

	p := &FakePresenter{}
	name := "name"
	itr := New(&FakeRepo{Names: []string{name}})
	itr.Execute(t.Context(), Request{Name: name, NewName: name}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	name := "name"
	itr := New(&FakeRepo{Names: []string{name},
		Err: e.ErrInternal})
	itr.Execute(t.Context(),
		Request{Name: name, NewName: name}, p)
	assert.True(t, p.GotInternalErr)
}
