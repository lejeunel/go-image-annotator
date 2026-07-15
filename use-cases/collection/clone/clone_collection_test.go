package clone

import (
	"context"
	"testing"

	task "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestSubmitTaskWithoutIdentity(t *testing.T) {
	itr := New(&FakeCloner{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.NotNil(t, p.GotErr)
	assert.False(t, p.GotSuccess)
}
func CreateCtxWithUserId(ctx context.Context, userId u.UserId) context.Context {
	user := u.NewUser(userId)
	return u.AppendUserToContext(ctx, user)
}
func TestHandleAuthErr(t *testing.T) {
	group := "my-group"
	itr := New(&FakeCloner{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"), Request{DestinationGroup: &group}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestRetrieveTaskSpecification(t *testing.T) {
	itr := New(&FakeCloner{})
	p := &FakePresenter{}
	userId := "user@mail.com"
	itr.Execute(CreateCtxWithUserId(t.Context(), userId), Request{}, p)
	assert.NotNil(t, p.Got)
	assert.Equal(t, task.CollectionCloneTask, p.Got.Type)
	assert.Equal(t, task.PendingTask, p.Got.State)
	assert.Equal(t, userId, p.Got.Issuer)
	assert.True(t, p.GotSuccess)
}
func TestHandleClonerErr(t *testing.T) {
	cloner := FakeCloner{Err: e.ErrInternal}
	itr := New(&cloner)
	p := &FakePresenter{}
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"), Request{}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
	assert.False(t, p.GotSuccess)
}

func TestSubmitTask(t *testing.T) {
	cloner := FakeCloner{}
	itr := New(&cloner)
	p := &FakePresenter{}
	userId := "user@mail.com"
	source := "source"
	destination := "destination"
	req := Request{Source: source, Destination: destination, Deep: true}
	itr.Execute(CreateCtxWithUserId(t.Context(), userId), req, p)
	assert.Equal(t, source, cloner.GotTask.Source)
	assert.Equal(t, destination, cloner.GotTask.Destination)
	assert.Equal(t, userId, cloner.GotTask.Issuer)
	assert.Equal(t, true, cloner.GotTask.Deep)
	assert.True(t, p.GotSuccess)
}
