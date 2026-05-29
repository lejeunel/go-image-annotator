package delete

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err                   error
	RemovedImageId        im.ImageId
	ErrOnRemoveImage      bool
	ErrOnRemoveAnnotation bool
}

func (r *FakeRepo) RemoveImageFromCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	if r.Err != nil {
		return r.Err
	}
	r.RemovedImageId = imageId
	return nil
}

func (r *FakeRepo) RemoveAnnotation(imageId im.ImageId, collectionId clc.CollectionId, annotationId a.AnnotationId) error {
	if r.ErrOnRemoveAnnotation {
		return r.Err
	}
	return nil
}

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(Response) {
	p.GotSuccess = true
}
