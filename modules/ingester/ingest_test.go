package ingester

import (
	"bytes"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingCollectionShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.CollectionRepo = &fk.CollectionRepo{ErrOnFind: e.ErrNotFound}
	_, err := ing.Ingest(Request{})
	assert.Error(t, err)
}

func TestHandleArtefactRepoError(t *testing.T) {
	ing := NewTestingIngester()
	ing.CollectionRepo = &fk.CollectionRepo{ErrOnFind: e.ErrInternal}
	_, err := ing.Ingest(Request{})
	assert.Error(t, err)
}

func TestNonExistingLabelShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.LabelRepo = &fk.LabelRepo{ErrOnFind: e.ErrNotFound}
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}})
	assert.Error(t, err)
}

func TestHandleLabelExistsInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.LabelRepo = &fk.LabelRepo{ErrOnFind: e.ErrInternal}
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}})
	assert.Error(t, err)
}

func TestHandleIngestionInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageRepo = &fk.ImageRepo{ErrOnAddToCollection: e.ErrInternal}
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.Error(t, err)
}

func TestHandleAddLabelInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.AnnotationRepo = &fk.AnnotationRepo{ErrOnAddLabel: e.ErrInternal}
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &fk.ImageReader{}})
	assert.Error(t, err)
}

func TestAddImageDuplicateHashShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageRepo = &fk.ImageRepo{HashAlreadyExists: true}
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrDuplicate)
}

func TestHandleDuplicateHashInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageRepo = &fk.ImageRepo{ErrOnFindHash: e.ErrInternal}
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestNonExistingBBoxLabelShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.LabelRepo = &fk.LabelRepo{ErrOnFind: e.ErrNotFound}
	_, err := ing.Ingest(Request{BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label"}},
		Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestHandleBoundingBoxValidationError(t *testing.T) {
	ing := NewTestingIngester()
	_, err := ing.Ingest(Request{BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: -2, Height: -4}},
		Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestHandleAddBoundingBoxInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.AnnotationRepo = &fk.AnnotationRepo{ErrOnAddBoundingBox: e.ErrInternal}
	_, err := ing.Ingest(Request{BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4}},
		Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnAddLabelMustDeleteImage(t *testing.T) {
	fileStore := &fk.FileStore{}
	imageRepo := &fk.ImageRepo{}
	ing := NewTestingIngester()
	ing.ArtefactRepo = fileStore
	ing.ImageRepo = imageRepo
	ing.AnnotationRepo = &fk.AnnotationRepo{ErrOnAddLabel: e.ErrInternal}
	ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &fk.ImageReader{}})
	assert.Equal(t, 1, imageRepo.NumDeletedImages)
	assert.Equal(t, 1, fileStore.NumDeletedImages)
}

func TestCorrectDataIsStored(t *testing.T) {
	artefactRepo := &fk.FileStore{}
	ing := NewTestingIngester()
	ing.ArtefactRepo = artefactRepo
	data := []byte("the-data")
	ing.Ingest(Request{Reader: &fk.ImageReader{Buffer: *bytes.NewBuffer(data)}})
	assert.True(t, bytes.Equal(artefactRepo.GotData, data))
}

func TestAddBoundingBoxToImage(t *testing.T) {
	annotationRepo := &fk.AnnotationRepo{}
	ing := NewTestingIngester()
	ing.AnnotationRepo = annotationRepo
	_, err := ing.Ingest(Request{BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4}},
		Reader: &fk.ImageReader{}})
	assert.NoError(t, err)
	assert.Equal(t, 1, annotationRepo.NumBoundingBoxesAdded)
}

func TestInternalErrOnAddImageShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageRepo = &fk.ImageRepo{ErrOnAddImage: e.ErrInternal}
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddImageWithHash(t *testing.T) {
	ing := NewTestingIngester()
	hash := []byte("the-hash")
	ing.Hasher = &fk.Hasher{Sum_: hash}
	imageRepo := &fk.ImageRepo{}
	ing.ImageRepo = imageRepo
	ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.True(t, bytes.Equal(imageRepo.GotHash, hash))
}

func TestAddImageLabel(t *testing.T) {
	ing := NewTestingIngester()
	annotationRepo := &fk.AnnotationRepo{}
	ing.AnnotationRepo = annotationRepo
	_, err := ing.Ingest(Request{Labels: []string{"a-label"},
		Reader: &fk.ImageReader{}})
	assert.NoError(t, err)
	assert.Equal(t, 1, annotationRepo.NumImageLabelsAdded)
}

func TestValidationErrOnImageMIMETypeInferShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageSpecsDetector = &fk.SpecsDetector{Err: e.ErrValidation}
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestShouldAddMIMEType(t *testing.T) {
	imageRepo := &fk.ImageRepo{}
	ing := NewTestingIngester()
	ing.ImageRepo = imageRepo
	specs := im.ImageSpecs{MIMEType: "image/jpeg"}
	ing.ImageSpecsDetector = &fk.SpecsDetector{Return: specs}
	ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.Equal(t, specs.MIMEType, imageRepo.GotSpecs.MIMEType)
}

func TestCollectionWithoutGroup(t *testing.T) {
	ing := NewTestingIngester()
	ing.CollectionRepo = &fk.CollectionRepo{Return: clc.NewCollection(clc.NewCollectionId(), "a-collection")}
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.NoError(t, err)
}

func TestShouldStoreIngestionTime(t *testing.T) {
	imageRepo := &fk.ImageRepo{}
	ing := NewTestingIngester()
	ing.ImageRepo = imageRepo
	specs := im.ImageSpecs{MIMEType: "image/jpeg"}
	now := time.Now()
	ing.clock = clockwork.NewFakeClockAt(now)
	ing.ImageSpecsDetector = &fk.SpecsDetector{Return: specs}
	ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.Equal(t, now, imageRepo.GotSpecs.IngestedAt)
}
