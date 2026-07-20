package clone

import (
	"context"
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestSubmitTaskWithoutIdentity(t *testing.T) {
	itr := NewTestingCloner()
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
	itr := NewTestingCloner()
	itr.Auth = fk.Auth{Err: e.ErrAuthorization}
	p := &FakePresenter{}
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"), Request{DestinationGroup: &group}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestReceiveTaskPayload(t *testing.T) {
	itr := NewTestingCloner()
	p := &FakePresenter{}
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{"source-collection"}}
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Source: "source-collection", Destination: "destination-collection"}, p)
	assert.Equal(t, task.CollectionCloneTask, p.Got.Type)
	assert.True(t, p.GotSuccess)
}

func TestCloningToAlreadyExistingCollectionShouldFail(t *testing.T) {
	itr := NewTestingCloner()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{"destination-collection"}}
	logger := &fk.EventLogger{}
	itr.EventLogger = logger
	p := &FakePresenter{}
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"), Request{Destination: "destination-collection"}, p)
	assert.Error(t, p.GotErr)
}

func TestErrorOnFindGroup(t *testing.T) {
	itr := NewTestingCloner()
	itr.GroupRepo = &fk.GroupRepo{ErrOnFind: e.ErrNotFound}
	logger := &fk.EventLogger{}
	itr.EventLogger = logger
	p := &FakePresenter{}
	dstGroup := "my-group"
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Destination: "destination-collection", DestinationGroup: &dstGroup}, p)
	assert.Error(t, p.GotErr)
}

func SetupCloneableCollection() (Interactor, clc.Collection, im.Image, *fk.ImageRepo, *fk.AnnotationRepo) {
	itr := NewTestingCloner()
	srcCollection := clc.NewCollection(clc.NewCollectionId(), "src")
	image := im.NewImage(im.NewImageId(), srcCollection)
	image.AddLabel(lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	imRepo := &fk.ImageRepo{IterateBaseImages: []im.BaseImage{{image.Id, srcCollection.Name}}}
	anRepo := &fk.AnnotationRepo{}
	itr.ImageRepo = imRepo
	itr.Store = &fk.ImageStore{Return: &image}
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{srcCollection.Name}}
	itr.AnnotationRepo = anRepo

	return itr, srcCollection, image, imRepo, anRepo

}

func TestCloneOneImage(t *testing.T) {
	itr, srcCollection, image, imRepo, _ := SetupCloneableCollection()
	p := &FakePresenter{}
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Source: srcCollection.Name, Destination: "destination-collection"}, p)
	assert.Equal(t, image.Id, imRepo.AddedImageId)
}

func TestDeepCloneAddsImageLabel(t *testing.T) {
	itr, srcCollection, image, _, annotationRepo := SetupCloneableCollection()
	p := &FakePresenter{}
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Source: srcCollection.Name, Destination: "destination-collection", Deep: true}, p)
	assert.Equal(t, image.Labels[0].Id, annotationRepo.AddedAnnotationId)
}
