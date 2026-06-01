package create

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(st.FakeProvider{}, Request{}, p)
	assert.Equal(t, true, p.GotAuthErr, "auth error")
	assert.Equal(t, false, p.GotSuccess)
}

func TestCreateCollectionWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Names: []string{name}})
	itr.Execute(st.FakeProvider{}, Request{Name: name}, p)
	assert.Equal(t, true, p.GotDuplicationErr, "duplication error")
	assert.Equal(t, false, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal},
		WithNameValidator(&v.FakeNameValidator{}))
	itr.Execute(st.FakeProvider{}, Request{}, p)
	assert.Equal(t, true, p.GotInternalErr)
}

func TestCreateCollectionWithInvalidNameShouldFail(t *testing.T) {
	name := "my-collection%/"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Names: []string{name}},
		WithNameValidator(&v.FakeNameValidator{Err: e.ErrValidation}))
	itr.Execute(st.FakeProvider{}, Request{Name: name}, p)
	assert.Equal(t, true, p.GotValidationErr, "name validatoin")
}

func TestCreateCollection(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	now := time.Now()
	itr := NewInteractor(repo, WithClock(clockwork.NewFakeClockAt(now)))
	req := Request{Name: "a-name", Description: "a-descriptin"}
	itr.Execute(st.FakeProvider{}, req, p)
	assert.Equal(t, repo.Got.Name, req.Name, "name")
	assert.Equal(t, repo.Got.Description, req.Description, "description")
	assert.Equal(t, repo.Got.CreatedAt, now, "creation date")
	assert.Equal(t, repo.Got.Id.IsNil(), false, "id")
}
