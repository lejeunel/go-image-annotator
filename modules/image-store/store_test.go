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
	ast "github.com/lejeunel/go-image-annotator/modules/file-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingCollectionShouldFail(t *testing.T) {
	s := New(&FakeRepo{MissingCollection: true}, &ast.FakeStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId().String(),
		Collection: "a-collection"})
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestErrOnFindLabelShouldFail(t *testing.T) {
	s := New(&FakeRepo{ErrOnFindImageLabel: true, Err: e.ErrInternal}, &ast.FakeStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId().String(),
		Collection: "a-collection"})
	assert.NotNil(t, err)
}

func TestErrOnFindBoundingBoxesShouldFail(t *testing.T) {
	s := New(&FakeRepo{ErrOnFindBoundingBoxes: true, Err: e.ErrInternal}, &ast.FakeStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId().String(),
		Collection: "a-collection"})
	assert.NotNil(t, err)
}

func TestErrOnExistsShouldFail(t *testing.T) {
	s := New(&FakeRepo{ErrOnExists: true, Err: e.ErrInternal}, &ast.FakeStore{})
	_, err := s.Find(im.BaseImage{ImageId: im.NewImageId().String(),
		Collection: "a-collection"})
	assert.NotNil(t, err)
}

func TestFindImageGivesCorrectAnnotations(t *testing.T) {
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	labels := []a.ImageLabel{{Id: a.NewAnnotationId(), Label: label}}
	bboxes := []a.BoundingBox{{Id: a.NewAnnotationId(), Label: label}}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")

	s := New(&FakeRepo{Collection: collection, Labels: labels,
		BoundingBoxes: bboxes}, &ast.FakeStore{Data: []byte("test-data")})
	image, _ := s.Find(im.BaseImage{ImageId: im.NewImageId().String(),
		Collection: collection.Name})
	assert.Equal(t, collection.Id, image.Collection.Id)
	assert.Equal(t, 1, len(image.Labels))
	assert.Equal(t, 1, len(image.BoundingBoxes))
}

func TestImageReaderGivesCorrectBytes(t *testing.T) {
	data := []byte("test-data")

	s := New(&FakeRepo{}, &ast.FakeStore{Data: data})
	image, _ := s.Find(im.BaseImage{ImageId: im.NewImageId().String(),
		Collection: "the-collection"})
	gotBytes, _ := io.ReadAll(image.Reader)
	assert.Equal(t, true, bytes.Equal(gotBytes, data))
}

func TestRetrieveSpecs(t *testing.T) {
	now := time.Now()
	s := New(&FakeRepo{Specs: im.ImageSpecs{IngestedAt: now}},
		&ast.FakeStore{})
	image, _ := s.Find(im.BaseImage{ImageId: im.NewImageId().String(),
		Collection: "the-collection"})
	assert.Equal(t, image.Specs.IngestedAt, now)
}
