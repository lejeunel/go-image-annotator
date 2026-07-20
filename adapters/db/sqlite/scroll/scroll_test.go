package scroll

import (
	"fmt"
	"testing"
	"time"

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
	assert.True(t, r.ImageId == ids[2])
}

func TestGettingPrevImage(t *testing.T) {
	repos := NewTestScrollerRepos()
	ids := CreateImagesWithOrderedIds(repos, 3)
	r, _ := repos.Scroller.GetAdjacent(ids[2], scr.NewCriteria(), scr.ScrollPrevious)
	assert.True(t, r.ImageId == ids[1])
}
func TestNextImageInCollection(t *testing.T) {
	repos := NewTestScrollerRepos()
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	repos.Collection.Create(collection)
	firstId, _ := im.NewImageIdFromString(FakeUUIDFromInt(0))
	secondId, _ := im.NewImageIdFromString(FakeUUIDFromInt(1))
	repos.Image.AddImage(firstId, []byte("first-hash"), im.ImageSpecs{})
	repos.Image.AddImage(secondId, []byte("second-hash"), im.ImageSpecs{})
	repos.Image.AddToCollection(firstId, collection.Id)
	repos.Image.AddToCollection(secondId, collection.Id)

	r, err := repos.Scroller.GetAdjacent(firstId,
		scr.NewCriteria(
			scr.WithCollection(collection.Name)),
		scr.ScrollNext)
	assert.NoError(t, err)
	assert.Equal(t, secondId, r.ImageId)
}

func CreateImagesWithIngestTime(repos SQLiteScrollerRepos, num int) ([]im.ImageId, clc.Collection) {
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	repos.Collection.Create(collection)
	ids := []im.ImageId{}
	now := time.Now()
	for n := range num {
		id := im.NewImageId()
		repos.Image.AddImage(id, fmt.Append([]byte{}, n),
			im.ImageSpecs{IngestedAt: now.Add(time.Duration(n) * time.Hour)})
		repos.Image.AddToCollection(id, collection.Id)
		ids = append(ids, id)
	}
	return ids, collection

}

func TestNextIngestedImageInCollection(t *testing.T) {
	repos := NewTestScrollerRepos()
	ids, collection := CreateImagesWithIngestTime(repos, 3)

	r, err := repos.Scroller.GetAdjacent(ids[0],
		scr.NewCriteria(
			scr.WithCollection(collection.Name),
			scr.WithOrdering(im.Ordering{IngestTime: true})),
		scr.ScrollNext)
	assert.NoError(t, err)
	assert.Equal(t, ids[1], r.ImageId)
}

func TestPreviousIngestedImageInCollection(t *testing.T) {
	repos := NewTestScrollerRepos()
	ids, collection := CreateImagesWithIngestTime(repos, 3)

	r, err := repos.Scroller.GetAdjacent(ids[2],
		scr.NewCriteria(
			scr.WithCollection(collection.Name),
			scr.WithOrdering(im.Ordering{IngestTime: true})),
		scr.ScrollPrevious)
	assert.NoError(t, err)
	assert.Equal(t, ids[1], r.ImageId)
}
