package scroll

import (
	"fmt"
	"testing"
	"time"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
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
	repos.Image.AddImage(id, nil, im.Specs{})
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
	id, _ := im.NewImageIdFromString(st.FakeUUIDFromInt(0))
	repos.Image.AddImage(id, nil, im.Specs{})
	_, err := repos.Scroller.GetAdjacent(id, scr.NewCriteria(), scr.ScrollPrevious)
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func CreateImagesWithIngestTime(repos SQLiteScrollerRepos, num int) ([]im.ImageId, clc.Collection) {
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	repos.Collection.Create(collection)
	ids := []im.ImageId{}
	now := time.Now()
	for n := range num {

		id, _ := im.NewImageIdFromString(st.FakeUUIDFromInt(n))
		repos.Image.AddImage(id, fmt.Append([]byte{}, n),
			im.Specs{IngestedAt: now.Add(time.Duration(n) * time.Hour)})
		repos.Image.AddToCollection(id, collection.Name)
		ids = append(ids, id)
	}
	return ids, collection
}

type ScrollerTest struct {
	name         string
	currentId    im.ImageId
	OrderByField string
	Order        im.Order
	Direction    scr.ScrollingDirection
	wantId       im.ImageId
}

func TestGetAdjacentImages(t *testing.T) {
	repos := NewTestScrollerRepos()
	ids, collection := CreateImagesWithIngestTime(repos, 3)

	tests := []ScrollerTest{
		{"next with ingested_at asc", ids[0], "ingested_at", im.AscOrder, scr.ScrollNext, ids[1]},
		{"prev with ingested_at asc", ids[2], "ingested_at", im.AscOrder, scr.ScrollPrevious, ids[1]},
		{"next with ingested_at desc", ids[2], "ingested_at", im.DescOrder, scr.ScrollNext, ids[1]},
		{"prev with ingested_at desc", ids[0], "ingested_at", im.DescOrder, scr.ScrollPrevious, ids[1]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := repos.Scroller.GetAdjacent(tt.currentId,
				scr.NewCriteria(
					scr.WithCollection(collection.Name),
					scr.WithOrdering(im.OrderingArgs{{Field: tt.OrderByField, Order: tt.Order}})),
				tt.Direction)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantId.String(), r.ImageId.String())

		})
	}
}
