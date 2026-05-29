package annotation

import (
	s "github.com/lejeunel/go-image-annotator/infra/db"
	sc "github.com/lejeunel/go-image-annotator/infra/db/collection"
	si "github.com/lejeunel/go-image-annotator/infra/db/image"
	sl "github.com/lejeunel/go-image-annotator/infra/db/label"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type AnnotationTestingRepos struct {
	Image      si.SQLiteImageRepo
	Collection sc.SQLiteCollectionRepo
	Label      sl.SQLiteLabelRepo
	Annotation SQLiteAnnotationRepo
}

func NewAnnotationTestRepos() AnnotationTestingRepos {
	db := s.NewSQLiteDB(":memory:")
	return AnnotationTestingRepos{Image: *si.NewSQLiteImageRepo(db),
		Collection: *sc.NewSQLiteCollectionRepo(db),
		Label:      *sl.NewSQLiteLabelRepo(db),
		Annotation: *NewSQLiteAnnotationRepo(db)}
}

func CreateAnnotableImage(repos AnnotationTestingRepos, collectionName string, labelName string) (im.Image, clc.Collection, lbl.Label) {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	label := lbl.NewLabel(lbl.NewLabelId(), labelName)
	repos.Label.Create(*label)
	repos.Collection.Create(*collection)
	image := im.NewImage(im.NewImageId(), *collection)
	repos.Image.AddImage(image.Id, nil, im.ImageSpecs{})
	repos.Image.AddToCollection(image.Id, collection.Id)

	return *image, *collection, *label

}
