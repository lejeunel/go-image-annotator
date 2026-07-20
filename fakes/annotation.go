package fake

import (
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type AnnotationRepo struct {
	Err                    error
	ErrOnAddPoly           error
	ErrOnAddLabel          error
	ErrOnAddBoundingBox    error
	ErrOnUpdate            error
	ErrOnFindPolygons      error
	ErrOnFindBoundingBoxes error
	ErrOnFindImageLabels   error
	ErrOnGetSpecs          error
	GotImageId             im.ImageId
	GotCollectionId        clc.CollectionId
	GotUserId              *u.UserId
	GotTime                *time.Time
	GotBox                 a.BoundingBox
	GotPolygon             a.Polygon
	AddedAnnotationId      a.AnnotationId
	AddedLabelId           lbl.LabelId
	AddedOnImageId         im.ImageId
	AddedOnCollectionId    clc.CollectionId
	GotUpdatableBox        a.BoundingBoxUpdatables
	GotUpdatablePoly       a.PolygonUpdatables
	GotRemovedAnnotation   a.AnnotationId
	ErrOnRemoveAnnotation  error
	UpdatedAnnotationId    a.AnnotationId
	UpdatedLabelId         lbl.LabelId
	NumBoundingBoxesAdded  int
	NumImageLabelsAdded    int
	Labels                 []a.ImageLabel
	BoundingBoxes          []a.BoundingBox
	Polygons               []a.Polygon
	Specs                  im.ImageSpecs

	NoGroup bool
}

func (r *AnnotationRepo) AddBoundingBox(imageId im.ImageId, collectionId clc.CollectionId, box a.BoundingBox, userId *u.UserId, t *time.Time) error {
	if r.ErrOnAddBoundingBox != nil {
		return r.ErrOnAddBoundingBox
	}
	r.GotImageId = imageId
	r.GotCollectionId = collectionId
	r.GotBox = box
	r.GotUserId = userId
	r.GotTime = t
	r.NumBoundingBoxesAdded += 1
	return nil
}

func (r *AnnotationRepo) AddPolygon(imageId im.ImageId, collectionId clc.CollectionId, poly a.Polygon, userId *u.UserId, t *time.Time) error {
	if r.ErrOnAddPoly != nil {
		return r.ErrOnAddPoly
	}
	r.GotImageId = imageId
	r.GotCollectionId = collectionId
	r.GotPolygon = poly
	r.GotUserId = userId
	r.GotTime = t
	return nil
}

func (r *AnnotationRepo) AddImageLabel(imageId im.ImageId, collectionId clc.CollectionId, imageLabel a.ImageLabel, userId *u.UserId, t *time.Time) error {
	if r.ErrOnAddLabel != nil {
		return r.ErrOnAddLabel
	}
	r.AddedAnnotationId = imageLabel.Id
	r.AddedLabelId = imageLabel.Label.Id
	r.AddedOnImageId = imageId
	r.AddedOnCollectionId = collectionId
	r.GotUserId = userId
	r.GotTime = t
	r.NumImageLabelsAdded += 1
	return nil

}

func (r *AnnotationRepo) UpdateBoundingBox(id a.AnnotationId, u a.BoundingBoxUpdatables, userId *u.UserId, t *time.Time) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdatableBox = u
	r.GotUserId = userId
	r.GotTime = t
	return nil
}

func (r *AnnotationRepo) GroupOfAnnotation(id a.AnnotationId) (*string, error) {
	if r.NoGroup {
		return nil, nil
	}
	group := "my-group"
	return &group, nil
}

func (r *AnnotationRepo) UpdatePolygon(id a.AnnotationId, u a.PolygonUpdatables, userId *u.UserId, t *time.Time) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdatablePoly = u
	r.GotUserId = userId
	r.GotTime = t
	return nil
}

func (r *AnnotationRepo) RemoveAnnotation(annotationId a.AnnotationId) error {
	if r.ErrOnRemoveAnnotation != nil {
		return r.ErrOnRemoveAnnotation
	}
	r.GotRemovedAnnotation = annotationId
	return nil
}

func (r *AnnotationRepo) UpdateLabelOfAnnotation(annotationId a.AnnotationId, labelId lbl.LabelId, userId *u.UserId, t *time.Time) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.UpdatedAnnotationId = annotationId
	r.UpdatedLabelId = labelId
	r.GotUserId = userId
	r.GotTime = t
	return nil
}

func (r *AnnotationRepo) FindBoundingBoxes(imageId im.ImageId, collectionId clc.CollectionId) ([]a.BoundingBox, error) {
	if r.ErrOnFindBoundingBoxes != nil {
		return nil, r.ErrOnFindBoundingBoxes
	}
	if r.BoundingBoxes != nil {
		return r.BoundingBoxes, nil
	}
	return nil, nil
}

func (r *AnnotationRepo) FindPolygons(imageId im.ImageId, collectionId clc.CollectionId) ([]a.Polygon, error) {
	if r.ErrOnFindPolygons != nil {
		return nil, r.ErrOnFindPolygons
	}
	if r.Polygons != nil {
		return r.Polygons, nil
	}
	return nil, nil
}

func (r *AnnotationRepo) FindImageLabels(imageId im.ImageId, collectionId clc.CollectionId) ([]a.ImageLabel, error) {
	if r.ErrOnFindImageLabels != nil {
		return nil, r.ErrOnFindImageLabels
	}
	if r.Labels != nil {
		return r.Labels, nil
	}
	return nil, nil
}

func (r *AnnotationRepo) GetSpecs(imageId im.ImageId) (*im.ImageSpecs, error) {
	if r.ErrOnGetSpecs != nil {
		return nil, r.Err
	}
	return &r.Specs, nil
}
