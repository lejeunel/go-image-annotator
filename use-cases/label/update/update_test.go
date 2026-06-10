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
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateNonExistingLabelShouldFail(t *testing.T) {

	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := New(&FakeRepo{})
	itr.Execute(t.Context(), Request{Name: non_existing_name}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateLabel(t *testing.T) {
	name := "name"

	p := &FakePresenter{}
	repo := &FakeRepo{Names: []string{name}}
	itr := New(repo)
	req := Request{Name: name,
		NewDescription: "updated-description"}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, p.Got.Description, req.NewDescription)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeErrRepo{e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
