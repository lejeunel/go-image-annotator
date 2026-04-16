package ingest

import (
	"bytes"
	"testing"

	ast "github.com/lejeunel/go-image-annotator-v2/application/file-store"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

func TestNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.CollectionRepo = &FakeCollectionRepo{MissingCollection: true}
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected not found error")
	}
}

func TestHandleInternalErrorOnCollectionExistsCheck(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.CollectionRepo = &FakeCollectionRepo{ErrOnFindCollection: true, Err: e.ErrInternal}
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestHandleArtefactRepoError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.CollectionRepo = &FakeCollectionRepo{ErrOnFindCollection: true, Err: e.ErrInternal}
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected invalid data error")
	}
}

func TestNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.LabelRepo = &FakeLabelRepo{MissingLabel: true}
	itr.Execute(Request{Labels: []string{"a-label"}}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected label not found error, but go none")
	}
}

func TestHandleLabelExistsInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.LabelRepo = &FakeLabelRepo{ErrOnLabelExists: true, Err: e.ErrInternal}
	itr.Execute(Request{Labels: []string{"a-label"}}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestHandleIngestionInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageRepo = &FakeImageRepo{ErrOnAddToCollection: true, Err: e.ErrInternal}
	itr.Execute(Request{Reader: &FakeImageReader{}}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestHandleAddLabelInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddLabel: true, Err: e.ErrInternal}
	itr.Execute(Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestHandleValidationErrorOnAddLabel(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.Execute(Request{Labels: []string{"a-label", "a-label"}, Reader: &FakeImageReader{}}, p)
	if !p.GotValidationErr || p.GotSuccess {
		t.Fatal("expected validation error")
	}
}

func TestAddImageDuplicateHashShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageRepo = &FakeImageRepo{HashAlreadyExists: true}
	itr.Execute(Request{Reader: &FakeImageReader{}}, p)
	if !p.GotDuplicationErr || p.GotSuccess {
		t.Fatal("expected validation error")
	}
}

func TestHandleDuplicateHashInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageRepo = &FakeImageRepo{ErrOnFindHash: true, Err: e.ErrInternal}
	itr.Execute(Request{Reader: &FakeImageReader{}}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestNonExistingBBoxLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.LabelRepo = &FakeLabelRepo{MissingLabel: true}
	itr.Execute(Request{BoundingBoxes: []BoundingBoxRequest{{Label: "a-label"}},
		Reader: &FakeImageReader{}}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected label not found error, but go none")
	}
}

func TestHandleBoundingBoxValidationError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.Execute(Request{BoundingBoxes: []BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: -2, Height: -4}},
		Reader: &FakeImageReader{}}, p)
	if !p.GotValidationErr || p.GotSuccess {
		t.Fatal("expected validation error")
	}
}

func TestHandleAddBoundingBoxInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddBoundingBox: true, Err: e.ErrInternal}
	itr.Execute(Request{BoundingBoxes: []BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4}},
		Reader: &FakeImageReader{}}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestInternalErrOnAddLabelMustDeleteImage(t *testing.T) {
	p := &FakePresenter{}
	fileStore := &ast.FakeStore{}
	imageRepo := &FakeImageRepo{}
	itr := NewTestingInteractor()
	itr.ArtefactRepo = fileStore
	itr.ImageRepo = imageRepo
	itr.AnnotationRepo = &FakeAnnotationRepo{ErrOnAddLabel: true, Err: e.ErrInternal}
	itr.Execute(Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}}, p)
	if imageRepo.NumDeletedImages != 1 || fileStore.NumDeletedImages != 1 {
		t.Fatalf("expected to delete image meta-data and artefacts, got %v and %v",
			imageRepo.NumDeletedImages, fileStore.NumDeletedImages)
	}
}

func TestCorrectDataIsStored(t *testing.T) {
	artefactRepo := &ast.FakeStore{}
	itr := NewTestingInteractor()
	itr.ArtefactRepo = artefactRepo
	data := []byte("the-data")
	itr.Execute(Request{Reader: &FakeImageReader{Buffer: *bytes.NewBuffer(data)}}, &FakePresenter{})
	if !bytes.Equal(artefactRepo.GotData, data) {
		t.Fatalf("stored bytes do no match input, got %v, input was %v",
			artefactRepo.GotData, data)
	}
}

func TestAddBoundingBoxToImage(t *testing.T) {
	p := &FakePresenter{}
	annotationRepo := &FakeAnnotationRepo{}
	itr := NewTestingInteractor()
	itr.AnnotationRepo = annotationRepo
	itr.Execute(Request{BoundingBoxes: []BoundingBoxRequest{{Label: "a-label", Xc: 10, Yc: 10, Width: 2, Height: 4}},
		Reader: &FakeImageReader{}}, p)
	if !p.GotSuccess || annotationRepo.NumBoundingboxesAdded != 1 {
		t.Fatalf("expected to add one bounding box to repo, but got %v", annotationRepo.NumBoundingboxesAdded)
	}
}

func TestInternalErrOnAddImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageRepo = &FakeImageRepo{ErrOnAddImage: true, Err: e.ErrInternal}
	itr.Execute(Request{Labels: []string{"a-label"}, Reader: &FakeImageReader{}}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestAddImageWithHash(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	hash := []byte("the-hash")
	itr.Hasher = &FakeHasher{sum: hash}
	imageRepo := &FakeImageRepo{}
	itr.ImageRepo = imageRepo
	itr.Execute(Request{Reader: &FakeImageReader{}}, p)
	if !p.GotSuccess || !bytes.Equal(imageRepo.GotHash, hash) {
		t.Fatalf("expected to store image hash %v, got %v", hash, imageRepo.GotHash)
	}
}

func TestAddImageWithTwoLabels(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	annotationRepo := &FakeAnnotationRepo{}
	itr.AnnotationRepo = annotationRepo
	itr.Execute(Request{Labels: []string{"a-label", "another-label"},
		Reader: &FakeImageReader{}}, p)
	if !p.GotSuccess || annotationRepo.NumLabelsAdded != 2 {
		t.Fatalf("expected to add two labels, but got %v", annotationRepo.NumLabelsAdded)
	}
}

func TestValidationErrOnImageMIMETypeInferShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ImageSpecsDetector = &FakeSpecsDetector{Err: e.ErrValidation}
	itr.Execute(Request{Reader: &FakeImageReader{}}, p)
	if !p.GotValidationErr || p.GotSuccess {
		t.Fatal("expected validation error")
	}
}

func TestAddImageShouldAddMIMEType(t *testing.T) {
	p := &FakePresenter{}
	imageRepo := &FakeImageRepo{}
	itr := NewTestingInteractor()
	itr.ImageRepo = imageRepo
	specs := im.ImageSpecs{MIMEType: "image/jpeg"}
	itr.ImageSpecsDetector = &FakeSpecsDetector{Return: specs}
	itr.Execute(Request{Reader: &FakeImageReader{}}, p)
	if imageRepo.GotSpecs.MIMEType != specs.MIMEType {
		t.Fatalf("expected to set MIMEType to %v, got %v", specs.MIMEType, imageRepo.GotSpecs.MIMEType)
	}
}

func TestAddImageWithMetaData(t *testing.T) {
}
