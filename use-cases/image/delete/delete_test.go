package delete

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := g.NewGroup(g.NewGroupId(), "my-group")
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection",
		clc.WithGroup(group.Name))
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.ImageStore{Return: &image}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: "a-collection"},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingResourceShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageStore{ErrOnFind: e.ErrNotFound})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestErrorOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageStore{ErrOnDelete: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDelete(t *testing.T) {
	p := &FakePresenter{}
	s := fk.ImageStore{}
	itr := New(&s)
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotSuccess)
	assert.NotNil(t, s.DeletedId)
}
