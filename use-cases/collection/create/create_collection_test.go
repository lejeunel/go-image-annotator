package create

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(context.Background(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateCollectionWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Names: []string{name}})
	itr.Execute(context.Background(), Request{Name: name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal},
		WithNameValidator(&v.FakeNameValidator{}))
	itr.Execute(context.Background(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}

func TestCreateCollectionWithInvalidNameShouldFail(t *testing.T) {
	name := "my-collection%/"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Names: []string{name}},
		WithNameValidator(&v.FakeNameValidator{Err: e.ErrValidation}))
	itr.Execute(context.Background(), Request{Name: name}, p)
	assert.True(t, p.GotValidationErr)
}

func TestCreateCollection(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	now := time.Now()
	itr := NewInteractor(repo, WithClock(clockwork.NewFakeClockAt(now)))
	req := Request{Name: "a-name", Description: "a-descriptin"}
	itr.Execute(context.Background(), req, p)
	assert.Equal(t, repo.Got.Name, req.Name)
	assert.Equal(t, repo.Got.Description, req.Description)
	assert.Equal(t, repo.Got.CreatedAt, now)
	assert.False(t, repo.Got.Id.IsNil())
}
