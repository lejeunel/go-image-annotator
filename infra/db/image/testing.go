package image

import (
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	s "github.com/lejeunel/go-image-annotator-v2/infra/db"
	cr "github.com/lejeunel/go-image-annotator-v2/infra/db/collection"
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
	repos.Image.AddImage(imageId, nil, "the-mimetype")

	return &imageId, &collectionId, repos.Image.AddToCollection(imageId, collectionId)
}
