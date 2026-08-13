package annotation

import (
	"github.com/jmoiron/sqlx"
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
	fk "github.com/lejeunel/go-image-annotator/fakes"
)

type AnnotationTestingRepos struct {
	Image      si.ImageRepo
	Collection sc.CollectionRepo
	Label      sl.LabelRepo
	Annotation AnnotationRepo
	Group      sg.GroupRepo
	User       su.UserRepo
}

func NewAnnotationTestRepos(db *sqlx.DB) AnnotationTestingRepos {
	return AnnotationTestingRepos{
		Image:      si.NewImageRepo(db, &fk.FilterStrParser{}, &fk.OrderStrParser{}),
		Collection: sc.NewCollectionRepo(db),
		Label:      sl.NewLabelRepo(db),
		Annotation: NewAnnotationRepo(db),
		Group:      sg.NewGroupRepo(db),
		User:       su.NewUserRepo(db),
	}
}

func CreateAnnotedImage(repos AnnotationTestingRepos, collectionName string, labelName string,
	group *string,
) (im.Image, clc.Collection, lbl.Label, a.ImageLabel) {
	image, collection, label := CreateAnnotableImage(repos, collectionName, labelName,
		group)
	imLabel := a.NewImageLabel(label)
	repos.Annotation.AddImageLabel(image.Id, collection.Name, imLabel, nil, nil)
	return image, collection, label, imLabel
}

func CreateAnnotableImage(repos AnnotationTestingRepos, collectionName string, labelName string,
	group *string,
) (im.Image, clc.Collection, lbl.Label) {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	if group != nil {
		group_ := grp.NewGroup(grp.NewGroupId(), *group)
		repos.Group.Create(group_)
		collection.Group = &group_.Name
	}
	label := lbl.NewLabel(lbl.NewLabelId(), labelName)
	repos.Label.Create(label)
	repos.Collection.Create(collection)
	image := im.NewImage(im.NewImageId(), collection)
	repos.Image.AddImage(image.Id, nil, im.Specs{})
	repos.Image.AddToCollection(image.Id, collection.Name)

	return image, collection, label
}
