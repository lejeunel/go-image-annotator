package query

import (
	sc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	si "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	sm "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/metadata"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	_ "modernc.org/sqlite"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"time"
)

func Setup() (sc.CollectionRepo, si.ImageRepo, sm.MetaRepo) {
	db := s.NewInMemory()
	filterParser, orderParser := si.MakeQueryParsers()
	imr := si.NewImageRepo(db, filterParser, orderParser)
	cr := sc.NewCollectionRepo(db)
	mr := sm.NewMetaRepo(db)
	return cr, imr, mr

}

func MustParseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

type TestIngestionPayload struct {
	ImageId       im.ImageId
	Collection    string
	IngestionTime time.Time
	MetaData      map[string]any
}
type FilterTest struct {
	name        string
	images      []TestIngestionPayload
	Filter      string
	Order       string
	WantFirstId im.ImageId
	WantCount   int
}

func findCollectionByName(cs []clc.Collection, name string) *clc.Collection {
	for _, c := range cs {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func InitFilterTest(imr si.ImageRepo, cr sc.CollectionRepo, sm sm.MetaRepo, payloads []TestIngestionPayload) {
	var createdCollections []clc.Collection
	for _, p := range payloads {
		collection := findCollectionByName(createdCollections, p.Collection)
		if collection == nil {
			c := clc.NewCollection(clc.NewCollectionId(), p.Collection)
			cr.Create(c)
			createdCollections = append(createdCollections, c)
			collection = &c
		}

		image := im.NewImage(p.ImageId, *collection)
		imr.AddImage(image.Id, []byte(image.Id.String()), im.Specs{IngestedAt: p.IngestionTime})
		imr.AddToCollection(image.Id, collection.Name)
		for k, v := range p.MetaData {
			sm.Add(collection.Name, image.Id, k, v)
		}
	}

}
