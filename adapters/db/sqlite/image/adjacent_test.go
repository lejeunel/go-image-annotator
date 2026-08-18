package image

import (
	"slices"
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
	_, err := repos.ImageRepo.GetAdjacent(im.NewImageId(), "", "", "", im.ScrollNext)
	assert.ErrorIs(t, err, e.ErrInternal)
}

type TestIngestionPayload struct {
	ImageId       im.ImageId
	Collection    string
	IngestionTime time.Time
}

type AdjacencyTest struct {
	name   string
	images []TestIngestionPayload
	im.FilterStr
	im.OrderStr
	currentImage       im.ImageId
	currentCollection  string
	wantPrev           *im.ImageId
	wantPrevCollection string
	wantNext           *im.ImageId
	wantNextCollection string
}

func findCollectionByName(cs []clc.Collection, name string) *clc.Collection {
	for _, c := range cs {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func InitAdjacencyTest(repos ScrollerRepos, payloads []TestIngestionPayload) {
	var createdCollections []clc.Collection
	var createdImages []im.ImageId
	for _, p := range payloads {
		collection := findCollectionByName(createdCollections, p.Collection)
		if collection == nil {
			c := clc.NewCollection(clc.NewCollectionId(), p.Collection)
			if err := repos.CollectionRepo.Create(c); err != nil {
				panic(err)
			}
			createdCollections = append(createdCollections, c)
			collection = &c
		}

		image := im.NewImage(p.ImageId, *collection)
		if !slices.Contains(createdImages, image.Id) {
			if err := repos.ImageRepo.AddImage(image.Id, []byte(image.Id.String()), im.Specs{IngestedAt: p.IngestionTime}); err != nil {
				panic(err)
			}
			createdImages = append(createdImages, image.Id)
		}
		if err := repos.ImageRepo.AddToCollection(image.Id, collection.Name); err != nil {
			panic(err)
		}
	}

}

var tests = []AdjacencyTest{
	{"single image has no adjacents",
		[]TestIngestionPayload{
			{*st.IdFromInt(0), "a-collection", time.Now()}},
		"",
		"",
		*st.IdFromInt(0),
		"a-collection",
		nil,
		"",
		nil,
		"",
	},
	{"one image per collection has no adjacents",
		[]TestIngestionPayload{
			{*st.IdFromInt(0), "a-collection", time.Now()},
			{*st.IdFromInt(1), "another-collection", time.Now()}},
		"collection:\"a-collection\"",
		"",
		*st.IdFromInt(0),
		"a-collection",
		nil,
		"",
		nil,
		"",
	},
	{"two images in one collection",
		[]TestIngestionPayload{
			{*st.IdFromInt(0), "a-collection", time.Now()},
			{*st.IdFromInt(1), "a-collection", time.Now()},
			{*st.IdFromInt(2), "another-collection", time.Now()}},
		"collection:\"a-collection\"",
		"ingested_at",
		*st.IdFromInt(0),
		"a-collection",
		nil,
		"",
		st.IdFromInt(1),
		"a-collection",
	},
	{"order by ingestion time ascending",
		[]TestIngestionPayload{
			{*st.IdFromInt(2), "a-collection", time.Now()},
			{*st.IdFromInt(1), "a-collection", time.Now()},
			{*st.IdFromInt(0), "a-collection", time.Now()},
		},
		"",
		"ingested_at",
		*st.IdFromInt(1),
		"a-collection",
		st.IdFromInt(2),
		"a-collection",
		st.IdFromInt(0),
		"a-collection",
	},
	{"order by ingestion time descending",
		[]TestIngestionPayload{
			{*st.IdFromInt(0), "a-collection", time.Now()},
			{*st.IdFromInt(1), "a-collection", time.Now()},
			{*st.IdFromInt(2), "a-collection", time.Now()},
		},
		"",
		"ingested_at:desc",
		*st.IdFromInt(1),
		"a-collection",
		st.IdFromInt(2),
		"a-collection",
		st.IdFromInt(0),
		"a-collection",
	},
	{"order by id by default",
		[]TestIngestionPayload{
			{*st.IdFromInt(1), "a-collection", time.Now()},
			{*st.IdFromInt(0), "a-collection", time.Now()},
			{*st.IdFromInt(2), "a-collection", time.Now()},
		},
		"",
		"",
		*st.IdFromInt(1),
		"a-collection",
		st.IdFromInt(0),
		"a-collection",
		st.IdFromInt(2),
		"a-collection",
	},
	{"order by id desc",
		[]TestIngestionPayload{
			{*st.IdFromInt(1), "a-collection", time.Now()},
			{*st.IdFromInt(2), "a-collection", time.Now()},
			{*st.IdFromInt(0), "a-collection", time.Now()},
		},
		"",
		"image_id:desc",
		*st.IdFromInt(1),
		"a-collection",
		st.IdFromInt(2),
		"a-collection",
		st.IdFromInt(0),
		"a-collection",
	},
}

func TestAdjacency(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repos := NewTestScrollerRepos(s.NewInMemory())
			InitAdjacencyTest(repos, tt.images)
			prev, err := repos.ImageRepo.GetAdjacent(tt.currentImage, tt.currentCollection, tt.FilterStr, tt.OrderStr, im.ScrollPrevious)
			assert.NoError(t, err)
			next, err := repos.ImageRepo.GetAdjacent(tt.currentImage, tt.currentCollection, tt.FilterStr, tt.OrderStr, im.ScrollNext)
			assert.NoError(t, err)
			if tt.wantPrev == nil {
				assert.Nil(t, prev, "previous")
			} else {
				assert.NotNil(t, prev)
				assert.Equal(t, tt.wantPrev.String(), prev.ImageId.String())
			}
			if tt.wantNext == nil {
				assert.Nil(t, next, "next")
			} else {
				assert.NotNil(t, next)
				assert.Equal(t, tt.wantNext.String(), next.ImageId.String())
			}
		})
	}
}
