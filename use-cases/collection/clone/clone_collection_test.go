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
	itr.Auth = fk.Auth{Err: e.ErrAuthorization}
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
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{"source-collection"}}
	itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Source: "source-collection", Destination: "destination-collection"}, p)
	assert.Equal(t, task.CollectionCloneTask, p.Got.Type)
	assert.True(t, p.GotSuccess)
}

func TestCloningToAlreadyExistingCollectionShouldFail(t *testing.T) {
	itr := NewTestingCloner()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{"destination-collection"}}
	p := &FakePresenter{}
	itr.Execute(
		st.CreateCtxWithUserId(t.Context(), "user@mail.com"),
		Request{Destination: "destination-collection"},
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
	collectionRepo := &fk.CollectionRepo{ExistingNames: []string{src.Name}, Return: src}
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
