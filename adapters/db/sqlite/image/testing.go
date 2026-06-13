package image

import (
	r "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	c "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

func MakeRepos() (SQLiteImageRepo, c.SQLiteCollectionRepo) {
	db := r.NewSQLiteDB(":memory:")
	return NewSQLiteImageRepo(db), c.NewSQLiteCollectionRepo(db)
}

func AddToCollection(imRepo SQLiteImageRepo, clcRepo c.SQLiteCollectionRepo,
	collectionName string, hash string) (*im.ImageId, *clc.CollectionId, error) {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	clcRepo.Create(collection)
	imageId := im.NewImageId()
	imRepo.AddImage(imageId, nil, im.ImageSpecs{})

	return &imageId, &collection.Id, imRepo.AddToCollection(imageId, collection.Id)
}
