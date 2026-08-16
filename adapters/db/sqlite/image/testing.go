package image

import (
	"github.com/jmoiron/sqlx"
	c "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	sc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
)

func BaseSetup() (ImageRepo, *sqlx.DB) {
	db := s.NewInMemory()
	return NewImageRepo(db, &fk.FilterStrParser{}, &fk.OrderStrParser{}), db
}

func SetupAdd(db *sqlx.DB) (ImageRepo, c.CollectionRepo, *sqlx.DB) {
	filterParser, orderParser := MakeQueryParsers()
	imr := NewImageRepo(db, filterParser, orderParser)
	return imr, c.NewCollectionRepo(db), db
}

func SetupList() (ImageRepo, sc.CollectionRepo, *sqlx.DB) {
	db := s.NewInMemory()
	imr, cr, _ := SetupAdd(db)
	return imr, cr, db
}

func AddToCollection(imRepo ImageRepo, clcRepo c.CollectionRepo,
	collectionName string, hash string,
) (*im.ImageId, *clc.CollectionId, error) {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	clcRepo.Create(collection)
	imageId := im.NewImageId()
	imRepo.AddImage(imageId, nil, im.Specs{})
	err := imRepo.AddToCollection(imageId, collection.Name)
	return &imageId, &collection.Id, err
}

type ImageListingTestingRepos struct {
	Image      ImageRepo
	Collection sc.CollectionRepo
}

func CreateSingleImageCollection(
	imr ImageRepo,
	cr c.CollectionRepo,
	collectionName string,
) (im.Image, clc.Collection) {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	cr.Create(collection)
	imageId := im.NewImageId()
	image := im.NewImage(imageId, collection)
	imr.AddImage(image.Id, nil, im.Specs{})
	imr.AddToCollection(image.Id, collection.Name)
	return image, collection
}

type ScrollerRepos struct {
	ImageRepo
	c.CollectionRepo
}

func NewTestScrollerRepos(db *sqlx.DB) ScrollerRepos {
	fp, op := MakeQueryParsers()
	return ScrollerRepos{
		ImageRepo:      NewImageRepo(db, fp, op),
		CollectionRepo: c.NewCollectionRepo(db),
	}
}

func CreateImagesWithOrderedIds(repos ScrollerRepos, num int) []im.ImageId {
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	repos.CollectionRepo.Create(collection)
	ids := []im.ImageId{}
	for n := range num {
		id, _ := im.NewImageIdFromString(st.FakeUUIDFromInt(n))
		repos.ImageRepo.AddImage(id, []byte(id.String()), im.Specs{})
		repos.ImageRepo.AddToCollection(id, collection.Name)
		ids = append(ids, id)
	}
	return ids
}

func CreateImageInCollection(repos ScrollerRepos,
	imageId im.ImageId, collectionName string,
) im.Image {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.CollectionRepo.Create(collection)
	image := im.NewImage(im.NewImageId(), collection)
	repos.ImageRepo.AddImage(image.Id, []byte(image.Id.String()), im.Specs{})
	repos.ImageRepo.AddToCollection(image.Id, image.Collection.Name)
	return image
}
