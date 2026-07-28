package delete

import (
	"context"
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func Setup(t *testing.T) (Interactor, clc.Collection, grp.Group, context.Context) {
	group := grp.NewGroup(grp.NewGroupId(), "my-group")
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection", clc.WithGroup(group))
	user := u.NewUser("user@mail.com", u.WithGroups([]string{"my-group"}))
	repos := Repos{
		CollectionRepo: &fk.CollectionRepo{
			Return: collection},
		ImageRepo:      &fk.ImageRepo{},
		AnnotationRepo: &fk.AnnotationRepo{},
	}
	itr := New(
		repos,
		&TestingTransactor{repos},
		&fk.ImageStore{},
		&fk.JobQueue{},
		&fk.EventLogger{}, fk.NewLogger(),
		WithAuth(fk.Auth{}))
	return itr, collection, group, u.AppendUserToContext(t.Context(), user)

}

func TestHandleAuthError(t *testing.T) {
	itr, _, _, _ := Setup(t)
	p := &FakePresenter{}
	itr.Auth = &fk.Auth{Err: e.ErrAuthorization}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingCollectionShouldFail(t *testing.T) {
	itr, _, _, _ := Setup(t)
	itr.CollectionRepo = &fk.CollectionRepo{ErrOnFind: e.ErrNotFound}
	p := &FakePresenter{}
	itr.Execute(t.Context(), "my-collection", p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnDelete(t *testing.T) {
	itr, _, _, _ := Setup(t)
	itr.CollectionRepo = &fk.CollectionRepo{ErrOnDelete: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(), "my-collection", p)
	assert.True(t, p.GotInternalErr)
}

func TestDeleteEmptyCollection(t *testing.T) {
	itr, _, _, ctx := Setup(t)
	p := &FakePresenter{}
	itr.Execute(ctx, "my-collection", p)
	assert.True(t, p.GotSuccess)
}

func TestRemoveImageFromCollection(t *testing.T) {
	itr, _, _, ctx := Setup(t)
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	imRepo := &fk.ImageRepo{IterateBaseImages: []im.BaseImage{{ImageId: image.Id, Collection: collection.Name}}}
	repos := Repos{
		CollectionRepo: &fk.CollectionRepo{
			Return: collection},
		ImageRepo:      imRepo,
		AnnotationRepo: &fk.AnnotationRepo{},
	}
	transactor := &TestingTransactor{repos}
	itr.Transactor = transactor
	itr.Repos = repos
	p := &FakePresenter{}
	itr.Execute(ctx, "my-collection", p)
	assert.Equal(t, image.Id, imRepo.RemovedImageId)
	assert.True(t, p.GotSuccess)
}

func TestRemoveAllAnnotations(t *testing.T) {
	itr, _, _, ctx := Setup(t)
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	imRepo := &fk.ImageRepo{IterateBaseImages: []im.BaseImage{{ImageId: image.Id, Collection: collection.Name}}}
	anRepo := &fk.AnnotationRepo{}
	repos := Repos{
		CollectionRepo: &fk.CollectionRepo{
			Return: collection},
		ImageRepo:      imRepo,
		AnnotationRepo: anRepo,
	}
	transactor := &TestingTransactor{repos}
	itr.Transactor = transactor
	itr.Repos = repos
	imageStore := &fk.ImageStore{}
	itr.ImageStore = imageStore
	p := &FakePresenter{}
	itr.Execute(ctx, "my-collection", p)
	assert.True(t, anRepo.RemovedAllAnnotations)
	assert.Equal(t, image.Id, imageStore.DeletedId)
	assert.True(t, p.GotSuccess)
}
