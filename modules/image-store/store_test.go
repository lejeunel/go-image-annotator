package image_store

import (
	"bytes"
	"io"
	"testing"
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingCollectionShouldFail(t *testing.T) {
	s := New(&fk.ImageRepo{}, &fk.CollectionRepo{}, &fk.AnnotationRepo{}, &fk.FileStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId(),
		Collection: "a-collection"})
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestErrOnFindLabelShouldFail(t *testing.T) {
	s := New(&fk.ImageRepo{}, &fk.CollectionRepo{}, &fk.AnnotationRepo{ErrOnFindImageLabels: e.ErrInternal}, &fk.FileStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId(),
		Collection: "a-collection"})
	assert.NotNil(t, err)
}

func TestErrOnFindBoundingBoxesShouldFail(t *testing.T) {
	s := New(&fk.ImageRepo{}, &fk.CollectionRepo{}, &fk.AnnotationRepo{ErrOnFindBoundingBoxes: e.ErrInternal}, &fk.FileStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId(),
		Collection: "a-collection"})
	assert.NotNil(t, err)
}

func TestErrOnFindPolygonsShouldFail(t *testing.T) {
	s := New(&fk.ImageRepo{}, &fk.CollectionRepo{}, &fk.AnnotationRepo{ErrOnFindPolygons: e.ErrInternal}, &fk.FileStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId(),
		Collection: "a-collection"})
	assert.NotNil(t, err)
}

func TestErrOnExistsShouldFail(t *testing.T) {
	s := New(&fk.ImageRepo{ErrOnImageExistsInCollection: e.ErrInternal}, &fk.CollectionRepo{}, &fk.AnnotationRepo{}, &fk.FileStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId(),
		Collection: "a-collection"})
	assert.NotNil(t, err)
}

func TestFindImageGivesCorrectAnnotations(t *testing.T) {
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	labels := []a.ImageLabel{{Id: a.NewAnnotationId(), Label: label}}
	bboxes := []a.BoundingBox{{Id: a.NewAnnotationId(), Label: label}}
	polygons := []a.Polygon{{Id: a.NewAnnotationId(), Label: label}}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")

	s := New(&fk.ImageRepo{ImageIsInCollection: true}, &fk.CollectionRepo{ExistingNames: []string{collection.Name},
		Return: collection}, &fk.AnnotationRepo{Labels: labels,
		BoundingBoxes: bboxes, Polygons: polygons}, &fk.FileStore{Data: []byte("test-data")})
	image, err := s.Find(im.BaseImage{ImageId: im.NewImageId(),
		Collection: collection.Name})
	assert.NoError(t, err)
	assert.Equal(t, collection.Id, image.Collection.Id)
	assert.Equal(t, 1, len(image.Labels))
	assert.Equal(t, 1, len(image.BoundingBoxes))
	assert.Equal(t, 1, len(image.Polygons))
}

func TestImageReaderGivesCorrectBytes(t *testing.T) {
	data := []byte("test-data")

	s := New(&fk.ImageRepo{ImageIsInCollection: true}, &fk.CollectionRepo{ExistingNames: []string{"the-collection"}},
		&fk.AnnotationRepo{}, &fk.FileStore{Data: data})
	image, _ := s.Find(im.BaseImage{ImageId: im.NewImageId(),
		Collection: "the-collection"})
	gotBytes, _ := io.ReadAll(image.Reader)
	assert.Equal(t, true, bytes.Equal(gotBytes, data))
}

func TestRetrieveSpecs(t *testing.T) {
	now := time.Now()
	s := New(&fk.ImageRepo{ImageIsInCollection: true, ReturnSpecs: im.ImageSpecs{IngestedAt: now}},
		&fk.CollectionRepo{ExistingNames: []string{"the-collection"}}, &fk.AnnotationRepo{},
		&fk.FileStore{})
	image, _ := s.Find(im.BaseImage{ImageId: im.NewImageId(),
		Collection: "the-collection"})
	assert.Equal(t, image.Specs.IngestedAt, now)
}
