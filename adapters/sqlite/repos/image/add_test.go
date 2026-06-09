package image

import (
	ist "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddSpecs(t *testing.T) {
	repos := NewImageTestRepos()
	id := im.NewImageId()

	specs := im.ImageSpecs{MIMEType: "the-mimetype", Width: 15, Height: 10}
	repos.Image.AddImage(id, nil, specs)
	r, err := repos.Image.GetSpecs(id)
	assert.NoError(t, err)
	assert.Equal(t, r.MIMEType, specs.MIMEType)
}

func TestCountAddedImageToCollection(t *testing.T) {
	repos := NewImageTestRepos()
	collection := "a-collection"
	AddToCollection(repos, collection, "")
	count, err := repos.Image.Count(ist.CountingParams{Collection: &collection})
	assert.NoError(t, err)
	assert.Equal(t, 1, int(*count))
}

func TestCountAllImagesWhenAddingImageToCollection(t *testing.T) {
	repos := NewImageTestRepos()
	AddToCollection(repos, "a-collection", "")
	count, err := repos.Image.Count(ist.CountingParams{})
	assert.NoError(t, err)
	assert.Equal(t, 1, int(*count))
}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	err := repo.AddToCollection(im.NewImageId(), clc.NewCollectionId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnIsCollectionPopulatedShouldFail(t *testing.T) {
	repos := NewImageTestRepos()
	collectionName := "a-collection"
	AddToCollection(repos, collectionName, "the-hash")
	repos.Image.Db.Close()
	_, err := repos.Collection.IsPopulated(collectionName)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestIsCollectionPopulated(t *testing.T) {
	repos := NewImageTestRepos()
	collectionName := "a-collection"
	AddToCollection(repos, collectionName, "the-hash")
	isPopulated, err := repos.Collection.IsPopulated(collectionName)
	assert.NoError(t, err)
	assert.True(t, *isPopulated)
}
