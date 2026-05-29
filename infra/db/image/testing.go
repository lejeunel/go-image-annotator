package image

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	s "github.com/lejeunel/go-image-annotator/infra/db"
	cr "github.com/lejeunel/go-image-annotator/infra/db/collection"
)

type ImageTestingRepos struct {
	Image      SQLiteImageRepo
	Collection cr.SQLiteCollectionRepo
}

func NewImageTestRepos() ImageTestingRepos {
	db := s.NewSQLiteDB(":memory:")
	return ImageTestingRepos{Image: *NewSQLiteImageRepo(db),
		Collection: *cr.NewSQLiteCollectionRepo(db)}
}

func AddToCollection(repos ImageTestingRepos, collectionName string, hash string) (*im.ImageId, *clc.CollectionId, error) {
	collectionId := clc.NewCollectionId()
	repos.Collection.Create(*clc.NewCollection(collectionId, collectionName))
	imageId := im.NewImageId()
	repos.Image.AddImage(imageId, nil, im.ImageSpecs{})

	return &imageId, &collectionId, repos.Image.AddToCollection(imageId, collectionId)
}
