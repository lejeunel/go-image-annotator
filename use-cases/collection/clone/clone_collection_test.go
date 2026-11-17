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
	itr, p := Setup()
	itr.Execute(t.Context(), Request{}, &p)
	assert.NotNil(t, p.GotErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleAuthErr(t *testing.T) {
	group := "my-group"
	itr, p := Setup()
	itr.Auth = fk.Auth{ErrOnAuth: e.ErrAuthorization}
	itr.Execute(
		st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{DestinationGroup: &group},
		&p,
	)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestReceiveTaskPayload(t *testing.T) {
	itr, p := Setup()
	sourceCollection := clc.NewCollection(clc.NewCollectionId(), "source-collection")
	destinationCollection := clc.NewCollection(clc.NewCollectionId(), "destination-collection")
	itr.CollectionRepo = &fk.CollectionRepo{Existing: []clc.Collection{sourceCollection}}
	itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{
			Source:      sourceCollection.Name,
			Destination: destinationCollection.Name,
		}, &p)
	assert.Equal(t, task.CollectionCloneTask.String(), p.Got.Type)
	assert.True(t, p.GotSuccess)
}

func TestCloningToAlreadyExistingCollectionShouldFail(t *testing.T) {
	itr, p := Setup()
	destinationCollection := clc.NewCollection(clc.NewCollectionId(), "destination-collection")
	itr.CollectionRepo = &fk.CollectionRepo{Existing: []clc.Collection{destinationCollection}}
	itr.Execute(
		st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Destination: destinationCollection.Name},
		&p,
	)
	assert.Error(t, p.GotErr)
}

func TestErrorOnFindGroup(t *testing.T) {
	itr, p := Setup()
	itr.GroupRepo = &fk.GroupRepo{ErrOnFind: e.ErrNotFound}
	dstGroup := "my-group"
	itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Destination: "destination-collection", DestinationGroup: &dstGroup}, &p)
	assert.Error(t, p.GotErr)
}

type cloneTest struct {
	itr            Interactor
	p              FakePresenter
	store          *fk.ImageStore
	collectionRepo *fk.CollectionRepo
	src            clc.Collection
	dst            string
}

func setupCloneTest() cloneTest {
	itr, p := Setup()
	src := clc.NewCollection(clc.NewCollectionId(), "source-collection",
		clc.WithCreatedAt(time.Now()), clc.WithGroup("a-group"), clc.WithProfile("a-profile"))
	store := &fk.ImageStore{}
	itr.ImageStore = store
	collectionRepo := &fk.CollectionRepo{Existing: []clc.Collection{src}}
	itr.CollectionRepo = collectionRepo
	itr.ImageRepo = &fk.ImageRepo{
		IterateBaseImages: []im.BaseImage{
			{ImageId: im.NewImageId(), Collection: src.Name},
			{ImageId: im.NewImageId(), Collection: src.Name},
		},
	}
	return cloneTest{itr, p, store, collectionRepo, src, "destination-collection"}
}

func TestCloneInheritsSourceAttributes(t *testing.T) {
	tc := setupCloneTest()
	tc.itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Source: tc.src.Name, Destination: tc.dst}, &tc.p)
	assert.Equal(t, *tc.src.Group, *tc.collectionRepo.Created.Group)
	assert.Equal(t, *tc.src.Profile, *tc.collectionRepo.Created.Profile)
	assert.NotEqual(t, tc.src.Id, tc.collectionRepo.Created.Id)
}

func TestShallowCloneByDefault(t *testing.T) {
	tc := setupCloneTest()
	tc.itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Source: tc.src.Name, Destination: tc.dst}, &tc.p)
	assert.Equal(t, tc.dst, tc.store.ShallowCopiedToCollection)
	assert.Empty(t, tc.store.DeepCopiedToCollection)
}

func TestDeepClone(t *testing.T) {
	tc := setupCloneTest()
	deep := true
	tc.itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Source: tc.src.Name, Destination: tc.dst, Deep: &deep}, &tc.p)
	assert.Equal(t, tc.dst, tc.store.DeepCopiedToCollection)
	assert.Empty(t, tc.store.ShallowCopiedToCollection)
}
