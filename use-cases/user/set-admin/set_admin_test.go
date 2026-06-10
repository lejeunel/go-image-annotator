package set_admin

import (
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user@example.com", true, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnSet(t *testing.T) {
	itr := New(&FakeRepo{Err: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user@example.com", true, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestSetAdmin(t *testing.T) {
	repo := &FakeRepo{}
	itr := New(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user@example.com", true, p)
	assert.True(t, p.GotSuccess)
	assert.True(t, p.Got.IsAdmin)
	assert.True(t, repo.GotValue)
}
