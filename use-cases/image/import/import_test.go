package import_image

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestNonExistingSourceImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ImageMissing: true})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected to get not found error")
	}
}

func TestInternalErrOnFindingSourceImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnImageExists: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected to get internal error")
	}
}

func TestNonExistingDestinationCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFindCollection: true, Err: e.ErrNotFound})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected to get not found error")
	}
}

func TestInternalErrOnFindCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFindCollection: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}

func TestImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ImageAlreadyInCollection: true})
	itr.Execute(Request{}, p)
	if !p.GotDependencyErr || p.GotSuccess {
		t.Fatalf("expected dependency error")
	}
}

func TestInternalErrOnImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnImageExistsInCollection: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}
func TestInternalErrOnImportShouldFail(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{ErrOnImport: true, Err: e.ErrInternal}
	itr := NewInteractor(repo)
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}

func TestImportImageInCollection(t *testing.T) {
	p := &FakePresenter{}
	imageId := im.NewImageId()
	collection := clc.NewCollection(clc.NewCollectionId(), "a-destination-collection")
	repo := &FakeRepo{DestinationCollection: *collection}
	itr := NewInteractor(repo)
	itr.Execute(Request{ImageId: imageId, Collection: collection.Name}, p)
	if !p.GotSuccess {
		t.Fatalf("expected success")
	}
	if (repo.ImportedImageId != imageId) || (repo.ImportedIntoCollectionId != collection.Id) {
		t.Fatalf("expected to import image %v into collection %v, got %v and %v",
			imageId, collection.Id, repo.ImportedImageId, repo.ImportedIntoCollectionId)
	}
}
