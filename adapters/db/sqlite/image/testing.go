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

func BaseSetup() (SQLiteImageRepo, *sqlx.DB) {
	db := s.NewInMemory()
	return NewSQLiteImageRepo(db, &fk.FilterStrParser{}, &fk.OrderStrParser{}), db
}

func SetupAdd(db *sqlx.DB) (SQLiteImageRepo, c.SQLiteCollectionRepo, *sqlx.DB) {
	filterParser, orderParser := MakeQueryParsers()
	imr := NewSQLiteImageRepo(db, filterParser, orderParser)
	return imr, c.NewSQLiteCollectionRepo(db), db
}

func SetupList() (SQLiteImageRepo, sc.SQLiteCollectionRepo, *sqlx.DB) {
	db := s.NewInMemory()
	imr, cr, _ := SetupAdd(db)
	return imr, cr, db

}

func AddToCollection(imRepo SQLiteImageRepo, clcRepo c.SQLiteCollectionRepo,
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
	Image      SQLiteImageRepo
	Collection sc.SQLiteCollectionRepo
}

func CreateSingleImageCollection(
	imr SQLiteImageRepo,
	cr c.SQLiteCollectionRepo,
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
