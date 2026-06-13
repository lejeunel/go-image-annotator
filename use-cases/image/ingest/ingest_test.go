package ingest

import (
	"bytes"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	ast "github.com/lejeunel/go-image-annotator/modules/file-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewTestingInteractor(WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.CollectionRepo = &FakeCollectionRepo{MissingCollection: true}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnCollectionExistsCheck(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.CollectionRepo = &FakeCollectionRepo{ErrOnFindCollection: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleArtefactRepoError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.CollectionRepo = &FakeCollectionRepo{ErrOnFindCollection: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.LabelRepo = &FakeLabelRepo{MissingLabel: true}
	itr.Execute(t.Context(), Request{Labels: []string{"a-label"}}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleLabelExistsInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.LabelRepo = &FakeLabelRepo{ErrOnLabelExists: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{Labels: []string{"a-label"}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleIngestionInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageRepo = &FakeImageRepo{ErrOnAddToCollection: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleAddLabelInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddLabel: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleValidationErrorOnAddLabel(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.Execute(t.Context(), Request{Labels: []string{"a-label", "a-label"}, Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestAddImageDuplicateHashShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageRepo = &FakeImageRepo{HashAlreadyExists: true}
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleDuplicateHashInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageRepo = &FakeImageRepo{ErrOnFindHash: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingBBoxLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.LabelRepo = &FakeLabelRepo{MissingLabel: true}
	itr.Execute(t.Context(), Request{BoundingBoxes: []BoundingBoxRequest{{Label: "a-label"}},
		Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleBoundingBoxValidationError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.Execute(t.Context(), Request{BoundingBoxes: []BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: -2, Height: -4}},
		Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleAddBoundingBoxInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddBoundingBox: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{BoundingBoxes: []BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4}},
		Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnAddLabelMustDeleteImage(t *testing.T) {
	p := &FakePresenter{}
	fileStore := &ast.FakeStore{}
	imageRepo := &FakeImageRepo{}
	itr := NewTestingInteractor()
	itr.ArtefactRepo = fileStore
	itr.ImageRepo = imageRepo
	itr.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddLabel: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}}, p)
	assert.Equal(t, 1, imageRepo.NumDeletedImages)
	assert.Equal(t, 1, fileStore.NumDeletedImages)
}

func TestCorrectDataIsStored(t *testing.T) {
	artefactRepo := &ast.FakeStore{}
	itr := NewTestingInteractor()
	itr.ArtefactRepo = artefactRepo
	data := []byte("the-data")
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{Buffer: *bytes.NewBuffer(data)}}, &FakePresenter{})
	assert.True(t, bytes.Equal(artefactRepo.GotData, data))
}

func TestAddBoundingBoxToImage(t *testing.T) {
	p := &FakePresenter{}
	annotationRepo := &FakeAnnotationRepo{}
	itr := NewTestingInteractor()
	itr.AnnotationRepo = annotationRepo
	itr.Execute(t.Context(), Request{BoundingBoxes: []BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4}},
		Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 1, annotationRepo.NumBoundingboxesAdded)
}

func TestInternalErrOnAddImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageRepo = &FakeImageRepo{ErrOnAddImage: true, Err: e.ErrInternal}
	itr.Execute(t.Context(), Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAddImageWithHash(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	hash := []byte("the-hash")
	itr.Hasher = &FakeHasher{sum: hash}
	imageRepo := &FakeImageRepo{}
	itr.ImageRepo = imageRepo
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{}}, p)
	assert.True(t, bytes.Equal(imageRepo.GotHash, hash))
	assert.True(t, p.GotSuccess)
}

func TestAddImageWithTwoLabels(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	annotationRepo := &FakeAnnotationRepo{}
	itr.AnnotationRepo = annotationRepo
	itr.Execute(t.Context(), Request{Labels: []string{"a-label", "another-label"},
		Reader: &FakeImageReader{}}, p)
	assert.Equal(t, 2, annotationRepo.NumLabelsAdded)
	assert.True(t, p.GotSuccess)
}

func TestValidationErrOnImageMIMETypeInferShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageSpecsDetector = &FakeSpecsDetector{Err: e.ErrValidation}
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestShouldAddMIMEType(t *testing.T) {
	imageRepo := &FakeImageRepo{}
	itr := NewTestingInteractor()
	itr.ImageRepo = imageRepo
	specs := im.ImageSpecs{MIMEType: "image/jpeg"}
	itr.ImageSpecsDetector = &FakeSpecsDetector{Return: specs}
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{}}, &FakePresenter{})
	assert.Equal(t, specs.MIMEType, imageRepo.GotSpecs.MIMEType)
}

func TestCollectionWithoutGroup(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.CollectionRepo = &FakeCollectionRepo{CollectionWithoutGroup: true}
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{}}, p)
	assert.True(t, p.GotSuccess)
}

func TestShouldStoreIngestionTime(t *testing.T) {
	imageRepo := &FakeImageRepo{}
	itr := NewTestingInteractor()
	itr.ImageRepo = imageRepo
	specs := im.ImageSpecs{MIMEType: "image/jpeg"}
	now := time.Now()
	itr.clock = clockwork.NewFakeClockAt(now)
	itr.ImageSpecsDetector = &FakeSpecsDetector{Return: specs}
	itr.Execute(t.Context(), Request{Reader: &FakeImageReader{}}, &FakePresenter{})
	assert.Equal(t, now, imageRepo.GotSpecs.IngestedAt)
}
