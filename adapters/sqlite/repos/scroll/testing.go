package scroll

import (
	"fmt"
	"strings"

	s "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	clcsql "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/collection"
	imsql "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/image"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type SQLiteScrollerRepos struct {
	Scroller   SQLiteScrollerRepo
	Image      imsql.SQLiteImageRepo
	Collection clcsql.SQLiteCollectionRepo
}

func NewTestScrollerRepos() *SQLiteScrollerRepos {
	db := s.NewSQLiteDB(":memory:")
	return &SQLiteScrollerRepos{
		Scroller:   NewSQLiteScrollerRepo(db),
		Image:      imsql.NewSQLiteImageRepo(db),
		Collection: clcsql.NewSQLiteCollectionRepo(db)}
}

func FakeUUIDFromInt(n int) string {
	digit := fmt.Sprintf("%d", n)
	full := strings.Repeat(digit, 32)
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		full[0:8],
		full[8:12],
		full[12:16],
		full[16:20],
		full[20:32],
	)
}

func CreateImagesWithOrderedIds(repo imsql.SQLiteImageRepo, num int) []im.ImageId {
	ids := []im.ImageId{}
	for n := range num {
		id, _ := im.NewImageIdFromString(FakeUUIDFromInt(n))
		repo.AddImage(id, []byte(id.String()), im.ImageSpecs{})
		ids = append(ids, id)
	}
	return ids
}

func CreateImageInCollection(imRepo imsql.SQLiteImageRepo, clcRepo clcsql.SQLiteCollectionRepo,
	imageId im.ImageId, collectionName string) im.Image {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	clcRepo.Create(collection)
	image := im.NewImage(im.NewImageId(), collection)
	imRepo.AddImage(image.Id, []byte(image.Id.String()), im.ImageSpecs{})
	imRepo.AddToCollection(image.Id, image.Collection.Id)
	return image
}
