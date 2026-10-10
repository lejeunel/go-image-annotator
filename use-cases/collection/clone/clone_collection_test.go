package clone

import (
	"testing"
	"time"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/entities/task"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	"github.com/stretchr/testify/assert"
)

func TestSubmitTaskWithoutIdentity(t *testing.T) {
	itr := NewTestingCloner()
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.NotNil(t, p.GotErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleAuthErr(t *testing.T) {
	group := "my-group"
	itr := NewTestingCloner()
	itr.Auth = fk.Auth{ErrOnAuth: e.ErrAuthorization}
	p := &FakePresenter{}
	itr.Execute(
		st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{DestinationGroup: &group},
		p,
	)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestReceiveTaskPayload(t *testing.T) {
	itr := NewTestingCloner()
	p := &FakePresenter{}
	sourceCollection := clc.NewCollection(clc.NewCollectionId(), "source-collection")
	destinationCollection := clc.NewCollection(clc.NewCollectionId(), "destination-collection")
	itr.CollectionRepo = &fk.CollectionRepo{Existing: []clc.Collection{sourceCollection}}
	itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{
			Source:      sourceCollection.Name,
			Destination: destinationCollection.Name,
		}, p)
	assert.Equal(t, task.CollectionCloneTask.String(), p.Got.Type)
	assert.True(t, p.GotSuccess)
}

func TestCloningToAlreadyExistingCollectionShouldFail(t *testing.T) {
	itr := NewTestingCloner()
	destinationCollection := clc.NewCollection(clc.NewCollectionId(), "destination-collection")
	itr.CollectionRepo = &fk.CollectionRepo{Existing: []clc.Collection{destinationCollection}}
	p := &FakePresenter{}
	itr.Execute(
		st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Destination: destinationCollection.Name},
		p,
	)
	assert.Error(t, p.GotErr)
}

func TestErrorOnFindGroup(t *testing.T) {
	itr := NewTestingCloner()
	itr.GroupRepo = &fk.GroupRepo{ErrOnFind: e.ErrNotFound}
	p := &FakePresenter{}
	dstGroup := "my-group"
	itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Destination: "destination-collection", DestinationGroup: &dstGroup}, p)
	assert.Error(t, p.GotErr)
}

func TestClone(t *testing.T) {
	itr := NewTestingCloner()
	dst := "destination-collection"
	src := clc.NewCollection(clc.NewCollectionId(), "source-collection",
		clc.WithCreatedAt(time.Now()), clc.WithGroup("a-group"), clc.WithProfile("a-profile"))
	s := fk.ImageStore{}
	itr.ImageStore = &s
	collectionRepo := &fk.CollectionRepo{Existing: []clc.Collection{src}}
	itr.CollectionRepo = collectionRepo
	itr.ImageRepo = &fk.ImageRepo{
		IterateBaseImages: []im.BaseImage{
			{ImageId: im.NewImageId(), Collection: src.Name},
			{ImageId: im.NewImageId(), Collection: src.Name},
		},
	}
	p := &FakePresenter{}
	itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Source: src.Name, Destination: dst}, p)
	assert.Equal(t, dst, s.CopiedToCollection)
	assert.Equal(t, *src.Group, *collectionRepo.Created.Group)
	assert.Equal(t, *src.Profile, *collectionRepo.Created.Profile)
	assert.NotEqual(t, src.Id, collectionRepo.Created.Id)
}
