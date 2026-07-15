package ingest

import (
	"bytes"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	ast "github.com/lejeunel/go-image-annotator/modules/file-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingCollectionShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.CollectionRepo = &FakeCollectionRepo{MissingCollection: true}
	_, err := ing.Ingest(Request{})
	assert.Error(t, err)
}

func TestHandleInternalErrorOnCollectionExistsCheck(t *testing.T) {
	ing := NewTestingIngester()
	ing.CollectionRepo = &FakeCollectionRepo{ErrOnFindCollection: true, Err: e.ErrInternal}
	_, err := ing.Ingest(Request{})
	assert.Error(t, err)
}

func TestHandleArtefactRepoError(t *testing.T) {
	ing := NewTestingIngester()
	ing.CollectionRepo = &FakeCollectionRepo{ErrOnFindCollection: true, Err: e.ErrInternal}
	_, err := ing.Ingest(Request{})
	assert.Error(t, err)
}

func TestNonExistingLabelShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.LabelRepo = &FakeLabelRepo{MissingLabel: true}
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}})
	assert.Error(t, err)
}

func TestHandleLabelExistsInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.LabelRepo = &FakeLabelRepo{ErrOnLabelExists: true, Err: e.ErrInternal}
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}})
	assert.Error(t, err)
}

func TestHandleIngestionInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageRepo = &FakeImageRepo{ErrOnAddToCollection: true, Err: e.ErrInternal}
	_, err := ing.Ingest(Request{Reader: &FakeImageReader{}})
	assert.Error(t, err)
}

func TestHandleAddLabelInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddLabel: true, Err: e.ErrInternal}
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}})
	assert.Error(t, err)
}

func TestHandleValidationErrorOnAddLabel(t *testing.T) {
	ing := NewTestingIngester()
	_, err := ing.Ingest(Request{Labels: []string{"a-label", "a-label"}, Reader: &FakeImageReader{}})
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestAddImageDuplicateHashShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageRepo = &FakeImageRepo{HashAlreadyExists: true}
	_, err := ing.Ingest(Request{Reader: &FakeImageReader{}})
	assert.ErrorIs(t, err, e.ErrDuplicate)
}

func TestHandleDuplicateHashInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageRepo = &FakeImageRepo{ErrOnFindHash: true, Err: e.ErrInternal}
	_, err := ing.Ingest(Request{Reader: &FakeImageReader{}})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestNonExistingBBoxLabelShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.LabelRepo = &FakeLabelRepo{MissingLabel: true}
	_, err := ing.Ingest(Request{BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label"}},
		Reader: &FakeImageReader{}})
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestHandleBoundingBoxValidationError(t *testing.T) {
	ing := NewTestingIngester()
	_, err := ing.Ingest(Request{BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: -2, Height: -4}},
		Reader: &FakeImageReader{}})
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestHandleAddBoundingBoxInternalErr(t *testing.T) {
	ing := NewTestingIngester()
	ing.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddBoundingBox: true, Err: e.ErrInternal}
	_, err := ing.Ingest(Request{BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4}},
		Reader: &FakeImageReader{}})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnAddLabelMustDeleteImage(t *testing.T) {
	fileStore := &ast.FakeStore{}
	imageRepo := &FakeImageRepo{}
	ing := NewTestingIngester()
	ing.ArtefactRepo = fileStore
	ing.ImageRepo = imageRepo
	ing.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddLabel: true, Err: e.ErrInternal}
	ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}})
	assert.Equal(t, 1, imageRepo.NumDeletedImages)
	assert.Equal(t, 1, fileStore.NumDeletedImages)
}

func TestCorrectDataIsStored(t *testing.T) {
	artefactRepo := &ast.FakeStore{}
	ing := NewTestingIngester()
	ing.ArtefactRepo = artefactRepo
	data := []byte("the-data")
	ing.Ingest(Request{Reader: &FakeImageReader{Buffer: *bytes.NewBuffer(data)}})
	assert.True(t, bytes.Equal(artefactRepo.GotData, data))
}

func TestAddBoundingBoxToImage(t *testing.T) {
	annotationRepo := &FakeAnnotationRepo{}
	ing := NewTestingIngester()
	ing.AnnotationRepo = annotationRepo
	_, err := ing.Ingest(Request{BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4}},
		Reader: &FakeImageReader{}})
	assert.NoError(t, err)
	assert.Equal(t, 1, annotationRepo.NumBoundingboxesAdded)
}

func TestInternalErrOnAddImageShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageRepo = &FakeImageRepo{ErrOnAddImage: true, Err: e.ErrInternal}
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddImageWithHash(t *testing.T) {
	ing := NewTestingIngester()
	hash := []byte("the-hash")
	ing.Hasher = &FakeHasher{sum: hash}
	imageRepo := &FakeImageRepo{}
	ing.ImageRepo = imageRepo
	ing.Ingest(Request{Reader: &FakeImageReader{}})
	assert.True(t, bytes.Equal(imageRepo.GotHash, hash))
}

func TestAddImageWithTwoLabels(t *testing.T) {
	ing := NewTestingIngester()
	annotationRepo := &FakeAnnotationRepo{}
	ing.AnnotationRepo = annotationRepo
	_, err := ing.Ingest(Request{Labels: []string{"a-label", "another-label"},
		Reader: &FakeImageReader{}})
	assert.Equal(t, 2, annotationRepo.NumLabelsAdded)
	assert.NoError(t, err)
}

func TestValidationErrOnImageMIMETypeInferShouldFail(t *testing.T) {
	ing := NewTestingIngester()
	ing.ImageSpecsDetector = &FakeSpecsDetector{Err: e.ErrValidation}
	_, err := ing.Ingest(Request{Reader: &FakeImageReader{}})
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestShouldAddMIMEType(t *testing.T) {
	imageRepo := &FakeImageRepo{}
	ing := NewTestingIngester()
	ing.ImageRepo = imageRepo
	specs := im.ImageSpecs{MIMEType: "image/jpeg"}
	ing.ImageSpecsDetector = &FakeSpecsDetector{Return: specs}
	ing.Ingest(Request{Reader: &FakeImageReader{}})
	assert.Equal(t, specs.MIMEType, imageRepo.GotSpecs.MIMEType)
}

func TestCollectionWithoutGroup(t *testing.T) {
	ing := NewTestingIngester()
	ing.CollectionRepo = &FakeCollectionRepo{CollectionWithoutGroup: true}
	_, err := ing.Ingest(Request{Reader: &FakeImageReader{}})
	assert.NoError(t, err)
}

func TestShouldStoreIngestionTime(t *testing.T) {
	imageRepo := &FakeImageRepo{}
	ing := NewTestingIngester()
	ing.ImageRepo = imageRepo
	specs := im.ImageSpecs{MIMEType: "image/jpeg"}
	now := time.Now()
	ing.clock = clockwork.NewFakeClockAt(now)
	ing.ImageSpecsDetector = &FakeSpecsDetector{Return: specs}
	ing.Ingest(Request{Reader: &FakeImageReader{}})
	assert.Equal(t, now, imageRepo.GotSpecs.IngestedAt)
}
