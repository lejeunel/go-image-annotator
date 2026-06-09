package create

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := "my-group"
	itr := NewInteractor(&FakeCollectionRepo{}, &FakeGroupRepo{Return: group}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Group: &group}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateCollectionWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeCollectionRepo{Names: []string{name}},
		&FakeGroupRepo{})
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeCollectionRepo{Err: e.ErrInternal},
		&FakeGroupRepo{},
		WithNameValidator(&v.FakeNameValidator{}))
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}

func TestCreateCollectionWithInvalidNameShouldFail(t *testing.T) {
	name := "my-collection%/"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeCollectionRepo{Names: []string{name}},
		&FakeGroupRepo{},
		WithNameValidator(&v.FakeNameValidator{Err: e.ErrValidation}))
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotValidationErr)
}

func TestCreateCollectionInNonExistingGroupShouldFail(t *testing.T) {
	name := "my-collection"
	group := "non-existing-group"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeCollectionRepo{},
		&FakeGroupRepo{MissingGroup: true},
	)
	itr.Execute(t.Context(), Request{Name: name, Group: &group}, p)
	assert.True(t, p.GotNotFoundErr)
}

func TestCreateCollection(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeCollectionRepo{}
	now := time.Now()
	group := "my-group"
	itr := NewInteractor(repo,
		&FakeGroupRepo{Return: group},
		WithClock(clockwork.NewFakeClockAt(now)))
	req := Request{Name: "a-name", Description: "a-description", Group: &group}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.Name, repo.Got.Name)
	assert.Equal(t, *req.Group, repo.Got.Group.Name)
	assert.Equal(t, req.Description, repo.Got.Description)
	assert.Equal(t, now, repo.Got.CreatedAt)
	assert.False(t, repo.Got.Id.IsNil())
}
