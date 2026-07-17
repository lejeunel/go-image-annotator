package create

import (
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := "my-group"
	itr := New(&fk.CollectionRepo{}, &fk.GroupRepo{Return: group}, WithAuth(fk.Auth{e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Group: &group}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateCollectionWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-collection"
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ErrOnCreate: e.ErrDuplicate},
		&fk.GroupRepo{})
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateCollectionWithInvalidNameShouldFail(t *testing.T) {
	name := "my-collection%/"
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ErrOnCreate: e.ErrValidation},
		&fk.GroupRepo{},
		WithNameValidator(&v.FakeNameValidator{Err: e.ErrValidation}))
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotValidationErr)
}

func TestCreateCollectionInNonExistingGroupShouldFail(t *testing.T) {
	name := "my-collection"
	group := "non-existing-group"
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{},
		&fk.GroupRepo{MissingGroup: true},
	)
	itr.Execute(t.Context(), Request{Name: name, Group: &group}, p)
	fmt.Println(p.GotErr)
	assert.True(t, p.GotNotFoundErr)
}

func TestCreateCollection(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.CollectionRepo{}
	now := time.Now()
	group := "my-group"
	itr := New(repo,
		&fk.GroupRepo{Return: group},
		WithClock(clockwork.NewFakeClockAt(now)))
	req := Request{Name: "a-name", Description: "a-description", Group: &group}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.Name, repo.Got.Name)
	assert.Equal(t, *req.Group, repo.Got.Group.Name)
	assert.Equal(t, req.Description, repo.Got.Description)
	assert.Equal(t, now, repo.Got.CreatedAt)
	assert.False(t, repo.Got.Id.IsNil())
}
