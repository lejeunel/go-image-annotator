package image

import (
	"github.com/jmoiron/sqlx"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	c "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	sc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
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
