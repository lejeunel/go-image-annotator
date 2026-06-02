package read

import (
	"github.com/stretchr/testify/assert"
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestReadCollection(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(),
		"my-collection",
		clc.WithDescription("a-description"))
	repo := &FakeRepo{Collection: collection}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	itr.Execute(t.Context(), Request{Name: collection.Name}, p)
	assert.Equal(t, Response{Name: collection.Name, Description: collection.Description}, p.Got)
}

func TestReadNonExistingCollectionShouldFail(t *testing.T) {
	repo := &FakeRepo{Collection: clc.Collection{Name: "my-collection", Description: "a-description"}}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	req := Request{Name: "non-existing-collection"}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}
