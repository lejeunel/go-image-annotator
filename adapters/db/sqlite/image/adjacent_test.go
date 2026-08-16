package image

import (
	"fmt"
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnImageMustExist(t *testing.T) {
	db := s.NewInMemory()
	repos := NewTestScrollerRepos(db)
	db.Close()
	_, err := repos.ImageRepo.ImageExists(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnGetAdjacent(t *testing.T) {
	db := s.NewInMemory()
	repos := NewTestScrollerRepos(db)
	db.Close()
	_, err := repos.ImageRepo.GetAdjacent(im.NewImageId(), "", "", im.ScrollNext)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestGettingAdjacentImageWhenSingle(t *testing.T) {
	repos := NewTestScrollerRepos(s.NewInMemory())
	id, _ := im.NewImageIdFromString(st.FakeUUIDFromInt(0))
	repos.ImageRepo.AddImage(id, nil, im.Specs{})
	rnext, err := repos.ImageRepo.GetAdjacent(id, "", "", im.ScrollNext)
	rprev, err := repos.ImageRepo.GetAdjacent(id, "", "", im.ScrollPrevious)
	assert.NoError(t, err)
	assert.Nil(t, rnext)
	assert.Nil(t, rprev)
}

func CreateImagesWithIngestTime(repos ScrollerRepos, num int) ([]im.ImageId, clc.Collection) {
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	repos.CollectionRepo.Create(collection)
	ids := []im.ImageId{}
	now := time.Now()
	for n := range num {
		id, _ := im.NewImageIdFromString(st.FakeUUIDFromInt(n))
		repos.ImageRepo.AddImage(id, fmt.Append([]byte{}, n),
			im.Specs{IngestedAt: now.Add(time.Duration(n) * time.Hour)})
		repos.ImageRepo.AddToCollection(id, collection.Name)
		ids = append(ids, id)
	}
	return ids, collection
}

type ScrollerTest struct {
	name      string
	currentId im.ImageId
	OrderBy   string
	Direction im.ScrollingDirection
	wantId    im.ImageId
}

func TestGetAdjacentImages(t *testing.T) {
	repos := NewTestScrollerRepos(s.NewInMemory())
	ids, collection := CreateImagesWithIngestTime(repos, 3)

	tests := []ScrollerTest{
		{"next with ingested_at asc", ids[0], "ingested_at", im.ScrollNext, ids[1]},
		{
			"prev with ingested_at asc",
			ids[2],
			"ingested_at:asc",
			im.ScrollPrevious,
			ids[1],
		},
		{"next with ingested_at desc", ids[2], "ingested_at:desc", im.ScrollNext, ids[1]},
		{
			"prev with ingested_at desc",
			ids[0],
			"ingested_at:desc",
			im.ScrollPrevious,
			ids[1],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := repos.ImageRepo.GetAdjacent(tt.currentId,
				fmt.Sprintf("collection:\"%v\"", collection.Name),
				tt.OrderBy,
				tt.Direction)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantId.String(), r.ImageId.String())
		})
	}
}
