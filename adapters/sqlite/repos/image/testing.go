package image

import (
	s "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	cr "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/collection"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
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
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(collection)
	imageId := im.NewImageId()
	repos.Image.AddImage(imageId, nil, im.ImageSpecs{})

	return &imageId, &collection.Id, repos.Image.AddToCollection(imageId, collection.Id)
}
