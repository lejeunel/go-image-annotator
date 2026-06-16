package annotation

import (
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	sc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	sg "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	si "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	sl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	su "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type AnnotationTestingRepos struct {
	Image      si.SQLiteImageRepo
	Collection sc.SQLiteCollectionRepo
	Label      sl.SQLiteLabelRepo
	Annotation SQLiteAnnotationRepo
	Group      sg.SQLiteGroupRepo
	User       su.SQLiteUserRepo
}

func NewAnnotationTestRepos() AnnotationTestingRepos {
	db := s.NewSQLiteDB(":memory:")
	return AnnotationTestingRepos{
		Image:      si.NewSQLiteImageRepo(db),
		Collection: sc.NewSQLiteCollectionRepo(db),
		Label:      sl.NewSQLiteLabelRepo(db),
		Annotation: NewSQLiteAnnotationRepo(db),
		Group:      sg.NewSQLiteGroupRepo(db),
		User:       su.NewSQLiteUserRepo(db)}
}
func CreateAnnotedImage(repos AnnotationTestingRepos, collectionName string, labelName string,
	group *string) (im.Image, clc.Collection, lbl.Label, a.ImageLabel) {
	image, collection, label := CreateAnnotableImage(repos, collectionName, labelName,
		group)
	imLabel := a.NewImageLabel(label)
	repos.Annotation.AddImageLabel(image.Id, collection.Id, imLabel, nil, nil)
	return image, collection, label, imLabel

}

func CreateAnnotableImage(repos AnnotationTestingRepos, collectionName string, labelName string,
	group *string) (im.Image, clc.Collection, lbl.Label) {

	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	if group != nil {
		group_ := grp.NewGroup(grp.NewGroupId(), *group)
		repos.Group.Create(group_)
		collection.Group = &group_
	}
	label := lbl.NewLabel(lbl.NewLabelId(), labelName)
	repos.Label.Create(label)
	repos.Collection.Create(collection)
	image := im.NewImage(im.NewImageId(), collection)
	repos.Image.AddImage(image.Id, nil, im.ImageSpecs{})
	repos.Image.AddToCollection(image.Id, collection.Id)

	return image, collection, label

}
