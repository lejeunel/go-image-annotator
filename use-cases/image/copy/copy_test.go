package copy

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

type copyTest struct {
	itr   Interactor
	p     FakePresenter
	store *fk.ImageStore
	req   Request
}

func setupCopyTest() copyTest {
	src := clc.NewCollection(clc.NewCollectionId(), "source")
	dst := clc.NewCollection(clc.NewCollectionId(), "destination", clc.WithGroup("a-group"))
	imageId := im.NewImageId()
	store := &fk.ImageStore{
		Existing: []im.BaseImage{{ImageId: imageId, Collection: src.Name}},
	}
	collections := &fk.CollectionRepo{Existing: []clc.Collection{src, dst}}
	return copyTest{
		itr:   New(store, collections),
		store: store,
		req: Request{
			ImageId:     imageId.String(),
			Source:      src.Name,
			Destination: dst.Name,
		},
	}
}

func (tc *copyTest) assertNothingCopied(t *testing.T) {
	assert.False(t, tc.p.GotSuccess)
	assert.Empty(t, tc.store.ShallowCopiedToCollection)
	assert.Empty(t, tc.store.DeepCopiedToCollection)
}

func TestCopyInvalidImageIdShouldFail(t *testing.T) {
	tc := setupCopyTest()
	tc.req.ImageId = "not-an-id"
	tc.itr.Execute(t.Context(), tc.req, &tc.p)
	assert.True(t, tc.p.GotValidationErr)
	tc.assertNothingCopied(t)
}

func TestCopyMissingSourceImageShouldFail(t *testing.T) {
	tc := setupCopyTest()
	tc.req.ImageId = im.NewImageId().String()
	tc.itr.Execute(t.Context(), tc.req, &tc.p)
	assert.True(t, tc.p.GotNotFoundErr)
	tc.assertNothingCopied(t)
}

func TestCopyMissingDestinationCollectionShouldFail(t *testing.T) {
	tc := setupCopyTest()
	tc.req.Destination = "does-not-exist"
	tc.itr.Execute(t.Context(), tc.req, &tc.p)
	assert.True(t, tc.p.GotNotFoundErr)
	tc.assertNothingCopied(t)
}

func TestCopyHandlesAuthError(t *testing.T) {
	tc := setupCopyTest()
	tc.itr.Auth = fk.Auth{ErrOnAuth: e.ErrAuthorization}
	tc.itr.Execute(t.Context(), tc.req, &tc.p)
	assert.True(t, tc.p.GotAuthErr)
	tc.assertNothingCopied(t)
}

func TestCopyImageAlreadyInDestinationShouldFail(t *testing.T) {
	tc := setupCopyTest()
	id, _ := im.NewImageIdFromString(tc.req.ImageId)
	tc.store.Existing = append(tc.store.Existing,
		im.BaseImage{ImageId: id, Collection: tc.req.Destination})
	tc.itr.Execute(t.Context(), tc.req, &tc.p)
	assert.True(t, tc.p.GotDuplicationErr)
	tc.assertNothingCopied(t)
}

func TestCopyIntoSameCollectionShouldFail(t *testing.T) {
	tc := setupCopyTest()
	tc.req.Destination = tc.req.Source
	tc.itr.Execute(t.Context(), tc.req, &tc.p)
	assert.True(t, tc.p.GotDuplicationErr)
	tc.assertNothingCopied(t)
}

func TestShallowCopy(t *testing.T) {
	tc := setupCopyTest()
	tc.itr.Execute(t.Context(), tc.req, &tc.p)
	assert.True(t, tc.p.GotSuccess)
	assert.Equal(t, tc.req.Destination, tc.store.ShallowCopiedToCollection)
	assert.Empty(t, tc.store.DeepCopiedToCollection)
	assert.Equal(t, tc.req.ImageId, tc.p.Got.Image.Id.String())
	assert.Equal(t, tc.req.Destination, tc.p.Got.Image.Collection.Name)
}

func TestDeepCopy(t *testing.T) {
	tc := setupCopyTest()
	tc.req.Deep = true
	tc.itr.Execute(t.Context(), tc.req, &tc.p)
	assert.True(t, tc.p.GotSuccess)
	assert.Equal(t, tc.req.Destination, tc.store.DeepCopiedToCollection)
	assert.Empty(t, tc.store.ShallowCopiedToCollection)
	assert.Equal(t, tc.req.Destination, tc.p.Got.Image.Collection.Name)
}
