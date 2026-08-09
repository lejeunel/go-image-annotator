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
	repos := NewTestingRepos()
	repos.CollectionRepo = &fk.CollectionRepo{ErrOnFind: e.ErrNotFound}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{})
	assert.Error(t, err)
}

func TestHandleArtefactRepoError(t *testing.T) {
	repos := NewTestingRepos()
	repos.CollectionRepo = &fk.CollectionRepo{ErrOnFind: e.ErrInternal}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{})
	assert.Error(t, err)
}

func TestNonExistingLabelShouldFail(t *testing.T) {
	repos := NewTestingRepos()
	repos.LabelRepo = &fk.LabelRepo{ErrOnFind: e.ErrNotFound}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}})
	assert.Error(t, err)
}

func TestHandleLabelExistsInternalErr(t *testing.T) {
	repos := NewTestingRepos()
	repos.LabelRepo = &fk.LabelRepo{ErrOnFind: e.ErrInternal}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}})
	assert.Error(t, err)
}

func TestHandleIngestionInternalErr(t *testing.T) {
	repos := NewTestingRepos()
	repos.ImageRepo = &fk.ImageRepo{ErrOnAddToCollection: e.ErrInternal}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.Error(t, err)
}

func TestHandleAddLabelInternalErr(t *testing.T) {
	repos := NewTestingRepos()
	repos.AnnotationRepo = &fk.AnnotationRepo{ErrOnAddLabel: e.ErrInternal}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &fk.ImageReader{}})
	assert.Error(t, err)
}

func TestAddImageDuplicateHashShouldFail(t *testing.T) {
	repos := NewTestingRepos()
	repos.ImageRepo = &fk.ImageRepo{HashAlreadyExists: true}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrDuplicate)
}

func TestHandleDuplicateHashInternalErr(t *testing.T) {
	repos := NewTestingRepos()
	repos.ImageRepo = &fk.ImageRepo{ErrOnFindHash: e.ErrInternal}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestNonExistingBBoxLabelShouldFail(t *testing.T) {
	repos := NewTestingRepos()
	repos.LabelRepo = &fk.LabelRepo{ErrOnFind: e.ErrNotFound}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{
		BoundingBoxes: []a.BoundingBoxRequest{{Label: "a-label"}},
		Reader:        &fk.ImageReader{},
	})
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestHandleBoundingBoxValidationError(t *testing.T) {
	repos := NewTestingRepos()
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{
		BoundingBoxes: []a.BoundingBoxRequest{
			{Label: "a-label", Xc: 10, Yc: 10, Width: -2, Height: -4},
		},
		Reader: &fk.ImageReader{},
	})
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestHandleAddBoundingBoxInternalErr(t *testing.T) {
	repos := NewTestingRepos()
	repos.AnnotationRepo = &fk.AnnotationRepo{ErrOnAddBoundingBox: e.ErrInternal}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{
		BoundingBoxes: []a.BoundingBoxRequest{
			{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4},
		},
		Reader: &fk.ImageReader{},
	})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnAddLabelMustDeleteImage(t *testing.T) {
	fileStore := &fk.FileStore{}
	repos := NewTestingRepos()
	repos.AnnotationRepo = &fk.AnnotationRepo{ErrOnAddLabel: e.ErrInternal}
	ing := NewTestingImageIngester(repos)
	ing.ArtefactRepo = fileStore
	ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &fk.ImageReader{}})
	assert.Equal(t, 1, fileStore.NumDeletedItems)
}

func TestCorrectDataIsStored(t *testing.T) {
	repos := NewTestingRepos()
	artefactRepo := &fk.FileStore{}
	ing := NewTestingImageIngester(repos)
	ing.ArtefactRepo = artefactRepo
	data := []byte("the-data")
	ing.Ingest(Request{Reader: &fk.ImageReader{Buffer: *bytes.NewBuffer(data)}})
	assert.True(t, bytes.Equal(artefactRepo.GotData, data))
}

func TestAddBoundingBoxToImage(t *testing.T) {
	repos := NewTestingRepos()
	anRepo := &fk.AnnotationRepo{}
	repos.AnnotationRepo = anRepo
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{
		BoundingBoxes: []a.BoundingBoxRequest{
			{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4},
		},
		Reader: &fk.ImageReader{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, anRepo.NumBoundingBoxesAdded)
}

func TestAddPolygon(t *testing.T) {
	repos := NewTestingRepos()
	anRepo := &fk.AnnotationRepo{}
	repos.AnnotationRepo = anRepo
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{
		Polygons: []a.PolygonRequest{
			{
				Label:  "a-label",
				Points: a.Points{Coordinates: [][2]float32{{0, 0}, {0, 1}}}},
		},
		Reader: &fk.ImageReader{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, anRepo.NumPolygonsAdded)
}

func TestInternalErrOnAddImageShouldFail(t *testing.T) {
	repos := NewTestingRepos()
	repos.ImageRepo = &fk.ImageRepo{ErrOnAddImage: e.ErrInternal}
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{Labels: []string{"a-label"}, Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddImageWithHash(t *testing.T) {
	repos := NewTestingRepos()
	imageRepo := &fk.ImageRepo{}
	repos.ImageRepo = imageRepo
	ing := NewTestingImageIngester(repos)
	hash := []byte("the-hash")
	ing.Hasher = &fk.Hasher{Sum_: hash}
	ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.True(t, bytes.Equal(imageRepo.GotHash, hash))
}

func TestAddImageLabel(t *testing.T) {
	repos := NewTestingRepos()
	annotationRepo := &fk.AnnotationRepo{}
	repos.AnnotationRepo = annotationRepo
	ing := NewTestingImageIngester(repos)
	_, err := ing.Ingest(Request{
		Labels: []string{"a-label"},
		Reader: &fk.ImageReader{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, annotationRepo.NumImageLabelsAdded)
}

func TestValidationErrOnImageMIMETypeInferShouldFail(t *testing.T) {
	repos := NewTestingRepos()
	ing := NewTestingImageIngester(repos)
	ing.ImageSpecsDetector = &fk.SpecsDetector{Err: e.ErrValidation}
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestShouldAddMIMEType(t *testing.T) {
	repos := NewTestingRepos()
	imageRepo := &fk.ImageRepo{}
	repos.ImageRepo = imageRepo
	ing := NewTestingImageIngester(repos)
	specs := im.Specs{MIMEType: "image/jpeg"}
	ing.ImageSpecsDetector = &fk.SpecsDetector{Return: specs}
	ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.Equal(t, specs.MIMEType, imageRepo.GotSpecs.MIMEType)
}

func TestCollectionWithoutGroup(t *testing.T) {
	repos := NewTestingRepos()
	ing := NewTestingImageIngester(repos)
	ing.CollectionRepo = &fk.CollectionRepo{
		Return: clc.NewCollection(clc.NewCollectionId(), "a-collection"),
	}
	_, err := ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.NoError(t, err)
}

func TestShouldStoreIngestionTime(t *testing.T) {
	repos := NewTestingRepos()
	imRepo := &fk.ImageRepo{}
	repos.ImageRepo = imRepo
	ing := NewTestingImageIngester(repos)
	specs := im.Specs{MIMEType: "image/jpeg"}
	now := time.Now()
	ing.Clock = clockwork.NewFakeClockAt(now)
	ing.ImageSpecsDetector = &fk.SpecsDetector{Return: specs}
	ing.Ingest(Request{Reader: &fk.ImageReader{}})
	assert.Equal(t, now, imRepo.GotSpecs.IngestedAt)
}

func TestHandleInvalidSpecsError(t *testing.T) {
	repos := NewTestingRepos()
	ing := NewTestingImageIngester(repos)
	ing.ImageSpecsDetector = &fk.SpecsDetector{Err: e.ErrValidation}
	_, err := ing.Ingest(Request{})
	assert.Error(t, err)
}
