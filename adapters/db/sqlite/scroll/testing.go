package scroll

import (
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	clcsql "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	imsql "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
)

type SQLiteScrollerRepos struct {
	Scroller   ScrollerRepo
	Image      imsql.ImageRepo
	Collection clcsql.CollectionRepo
}

func NewTestScrollerRepos() SQLiteScrollerRepos {
	db := s.NewInMemory()
	return SQLiteScrollerRepos{
		Scroller:   NewScrollerRepo(db),
		Image:      imsql.NewImageRepo(db, &fk.FilterStrParser{}, &fk.OrderStrParser{}),
		Collection: clcsql.NewCollectionRepo(db),
	}
}

func CreateImagesWithOrderedIds(repos SQLiteScrollerRepos, num int) []im.ImageId {
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	repos.Collection.Create(collection)
	ids := []im.ImageId{}
	for n := range num {
		id, _ := im.NewImageIdFromString(st.FakeUUIDFromInt(n))
		repos.Image.AddImage(id, []byte(id.String()), im.Specs{})
		repos.Image.AddToCollection(id, collection.Name)
		ids = append(ids, id)
	}
	return ids
}

func CreateImageInCollection(imRepo imsql.ImageRepo, clcRepo clcsql.CollectionRepo,
	imageId im.ImageId, collectionName string,
) im.Image {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	clcRepo.Create(collection)
	image := im.NewImage(im.NewImageId(), collection)
	imRepo.AddImage(image.Id, []byte(image.Id.String()), im.Specs{})
	imRepo.AddToCollection(image.Id, image.Collection.Name)
	return image
}
