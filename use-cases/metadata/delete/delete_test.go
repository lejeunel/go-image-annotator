package delete

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	group := g.NewGroup(g.NewGroupId(), "my-group")
	image.Collection.Group = &group.Name
	itr := New(&fk.CollectionRepo{ReturnGroup: group.Name}, &fk.MetaDataRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestErrorOnKeyExistsShouldFail(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{ErrOnKeyExists: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name, Key: "the-key"},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
	assert.False(t, p.GotSuccess)
}

func TestFailOnNonExistingKey(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name, Key: "the-key"},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
	assert.False(t, p.GotSuccess)
}

func TestErrorOnDelete(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	key := "the-key"
	itr := New(&fk.CollectionRepo{},
		&fk.MetaDataRepo{ExistingKeys: []string{key},
			ErrOnDelete: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name, Key: key},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
	assert.False(t, p.GotSuccess)

}

func TestDelete(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	key := "the-key"
	m := &fk.MetaDataRepo{ExistingKeys: []string{key}}
	itr := New(&fk.CollectionRepo{}, m)
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name, Key: key},
		p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, key, m.DeletedKey)
}
