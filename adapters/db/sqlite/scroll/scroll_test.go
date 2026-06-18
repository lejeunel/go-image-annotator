package scroll

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnImageMustExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	repos.Scroller.Db.Close()
	err := repos.Scroller.ImageMustExist(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnCollectionMustExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	repos.Scroller.Db.Close()
	err := repos.Scroller.CollectionMustExist("a-collection")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnGetAdjacent(t *testing.T) {
	repos := NewTestScrollerRepos()
	repos.Scroller.Db.Close()
	_, err := repos.Scroller.GetAdjacent(im.NewImageId(), scr.NewCriteria(), scr.ScrollNext)
	assert.ErrorIs(t, err, e.ErrInternal)

}

func TestShouldFailWhenImageDoesNotExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	err := repos.Scroller.ImageMustExist(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestImageMustExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	id := im.NewImageId()
	repos.Image.AddImage(id, nil, im.ImageSpecs{})
	err := repos.Scroller.ImageMustExist(id)
	assert.NoError(t, err)
}

func TestShouldFailWhenCollectionDoesNotExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	err := repos.Scroller.CollectionMustExist("non-existing-collection")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestShouldFailWhenNoImage(t *testing.T) {
	repos := NewTestScrollerRepos()
	id := im.NewImageId()
	_, err := repos.Scroller.GetAdjacent(id, scr.NewCriteria(), scr.ScrollNext)
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestGettingAdjacentImageWhenSingleImageShouldFail(t *testing.T) {
	repos := NewTestScrollerRepos()
	id, _ := im.NewImageIdFromString("00000000-0000-0000-0000-000000000000")
	repos.Image.AddImage(id, nil, im.ImageSpecs{})
	_, err := repos.Scroller.GetAdjacent(id, scr.NewCriteria(), scr.ScrollPrevious)
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestGettingNextImage(t *testing.T) {
	repos := NewTestScrollerRepos()
	ids := CreateImagesWithOrderedIds(repos, 3)
	r, err := repos.Scroller.GetAdjacent(ids[1], scr.NewCriteria(), scr.ScrollNext)
	assert.NoError(t, err)
	assert.True(t, r.ImageId == ids[2].String())
}

func TestGettingPrevImage(t *testing.T) {
	repos := NewTestScrollerRepos()
	ids := CreateImagesWithOrderedIds(repos, 3)
	r, _ := repos.Scroller.GetAdjacent(ids[2], scr.NewCriteria(), scr.ScrollPrevious)
	assert.True(t, r.ImageId == ids[1].String())
}

func TestScrollWithCollectionCriteria(t *testing.T) {
	repos := NewTestScrollerRepos()
	firstImage := CreateImageInCollection(repos.Image, repos.Collection,
		im.NewImageId(), "first-collection")
	CreateImageInCollection(repos.Image, repos.Collection,
		im.NewImageId(), "second-collection")

	_, err := repos.Scroller.GetAdjacent(firstImage.Id,
		scr.NewCriteria(scr.WithCollection("first-collection")),
		scr.ScrollPrevious)

	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestGettingNextImageInCollection(t *testing.T) {
	repos := NewTestScrollerRepos()
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	ids := CreateImagesWithOrderedIds(repos, 2)
	repos.Collection.Create(collection)
	repos.Image.AddToCollection(ids[0], collection.Id)
	repos.Image.AddToCollection(ids[1], collection.Id)

	r, _ := repos.Scroller.GetAdjacent(ids[0],
		scr.NewCriteria(scr.WithCollection(collection.Name)),
		scr.ScrollNext)
	assert.True(t, r.Collection == collection.Name)
}
