package query

import (
	"slices"
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
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

type AdjTestPayload struct {
	ImageId       im.ImageId
	Collection    string
	IngestionTime time.Time
	Meta          []m.MetaData
}

type AdjTest struct {
	name   string
	images []AdjTestPayload
	im.FilterStr
	im.OrderStr
	currentImage       im.ImageId
	currentCollection  string
	wantPrev           *im.ImageId
	wantPrevCollection string
	wantNext           *im.ImageId
	wantNextCollection string
}

func InitAdjacencyTest(repos ScrollerRepos, payloads []AdjTestPayload) {
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
		for _, m := range p.Meta {
			if err := repos.MetaRepo.Add(collection.Name, image.Id, m.Key, m.Value); err != nil {
				panic(err)
			}
		}
	}

}

var adjTests = []AdjTest{
	{"single image has no adjacents",
		[]AdjTestPayload{
			{*st.IdFromInt(0), "a-collection", time.Now(), nil}},
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
		[]AdjTestPayload{
			{*st.IdFromInt(0), "a-collection", time.Now(), nil},
			{*st.IdFromInt(1), "another-collection", time.Now(), nil}},
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
		[]AdjTestPayload{
			{*st.IdFromInt(0), "a-collection", time.Now(), nil},
			{*st.IdFromInt(1), "a-collection", time.Now(), nil},
			{*st.IdFromInt(2), "another-collection", time.Now(), nil}},
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
		[]AdjTestPayload{
			{*st.IdFromInt(2), "a-collection", time.Now(), nil},
			{*st.IdFromInt(1), "a-collection", time.Now(), nil},
			{*st.IdFromInt(0), "a-collection", time.Now(), nil},
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
		[]AdjTestPayload{
			{*st.IdFromInt(0), "a-collection", time.Now(), nil},
			{*st.IdFromInt(1), "a-collection", time.Now(), nil},
			{*st.IdFromInt(2), "a-collection", time.Now(), nil},
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
	{"order by image id by default",
		[]AdjTestPayload{
			{*st.IdFromInt(1), "a-collection", time.Now(), nil},
			{*st.IdFromInt(0), "a-collection", time.Now(), nil},
			{*st.IdFromInt(2), "a-collection", time.Now(), nil},
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
		[]AdjTestPayload{
			{*st.IdFromInt(1), "a-collection", time.Now(), nil},
			{*st.IdFromInt(2), "a-collection", time.Now(), nil},
			{*st.IdFromInt(0), "a-collection", time.Now(), nil},
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
	{"same image in two collections",
		[]AdjTestPayload{
			{*st.IdFromInt(0), "collection-v1", time.Now(), nil},
			{*st.IdFromInt(0), "collection-v2", time.Now(), nil}},
		"",
		"collection",
		*st.IdFromInt(0),
		"collection-v1",
		nil,
		"",
		st.IdFromInt(0),
		"collection-v2",
	},
	{"filter by meta",
		[]AdjTestPayload{
			{*st.IdFromInt(0), "a-collection", time.Now(), nil},
			{*st.IdFromInt(1), "a-collection", time.Now(), []m.MetaData{{Key: "is-active", Value: true}}},
			{*st.IdFromInt(2), "another-collection", time.Now(), []m.MetaData{{Key: "is-active", Value: true}}}},
		"meta.is-active?",
		"",
		*st.IdFromInt(1),
		"a-collection",
		nil,
		"",
		st.IdFromInt(2),
		"another-collection",
	},
}

func TestAdjacency(t *testing.T) {
	for _, tt := range adjTests {
		t.Run(tt.name, func(t *testing.T) {
			repos := NewTestScrollerRepos(s.NewInMemory())
			InitAdjacencyTest(repos, tt.images)

			adj, err := repos.ImageRepo.GetAdjacent(tt.currentImage, tt.currentCollection, tt.FilterStr, tt.OrderStr, im.ScrollPrevious)
			assert.NoError(t, err)
			if tt.wantPrev == nil {
				assert.Nil(t, adj.Prev, "previous")
			} else {
				assert.NotNil(t, adj.Prev)
				assert.Equal(t, tt.wantPrev.String(), adj.Prev.ImageId.String())
			}
			if tt.wantNext == nil {
				assert.Nil(t, adj.Next, "next")
			} else {
				assert.NotNil(t, adj.Next)
				assert.Equal(t, tt.wantNext.String(), adj.Next.ImageId.String())
			}
		})
	}
}
