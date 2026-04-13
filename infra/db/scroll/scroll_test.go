package scroll

import (
	"errors"
	"testing"

	"github.com/lejeunel/go-image-annotator-v2/application/scroller"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

func TestInternalErrOnImageMustExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	repos.Scroller.Db.Close()
	err := repos.Scroller.ImageMustExist(im.NewImageId())
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestInternalErrOnCollectionMustExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	repos.Scroller.Db.Close()
	err := repos.Scroller.CollectionMustExist("a-collection")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestInternalErrOnGetAdjacent(t *testing.T) {
	repos := NewTestScrollerRepos()
	repos.Scroller.Db.Close()
	_, err := repos.Scroller.GetAdjacent(im.NewImageId(), scroller.NewCriteria(), scroller.ScrollNext)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestShouldFailWhenImageDoesNotExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	err := repos.Scroller.ImageMustExist(im.NewImageId())
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error got %v", err)
	}
}

func TestImageMustExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	id := im.NewImageId()
	repos.Image.AddImage(id, nil, "")
	err := repos.Scroller.ImageMustExist(id)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
}

func TestShouldFailWhenCollectionDoesNotExist(t *testing.T) {
	repos := NewTestScrollerRepos()
	err := repos.Scroller.CollectionMustExist("non-existing-collection")
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error got %v", err)
	}
}

func TestShouldFailWhenNoImage(t *testing.T) {
	repos := NewTestScrollerRepos()
	id := im.NewImageId()
	_, err := repos.Scroller.GetAdjacent(id, scroller.NewCriteria(), scroller.ScrollNext)
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error got %v", err)
	}
}

func TestGettingAdjacentImageWhenSingleImageShouldFail(t *testing.T) {
	repos := NewTestScrollerRepos()
	id, _ := im.NewImageIdFromString("00000000-0000-0000-0000-000000000000")
	repos.Image.AddImage(id, nil, "")
	_, err := repos.Scroller.GetAdjacent(id, scroller.NewCriteria(), scroller.ScrollPrevious)
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found err got  %v", err)
	}
}

func TestGettingNextImage(t *testing.T) {
	repos := NewTestScrollerRepos()
	ids := CreateImagesWithOrderedIds(repos.Image, 3)
	r, _ := repos.Scroller.GetAdjacent(ids[1], scroller.NewCriteria(), scroller.ScrollNext)
	if r.ImageId != ids[2] {
		t.Fatalf("expected to retrieve next image with id %v got %v", ids[2], r.ImageId)
	}
}

func TestGettingPrevImage(t *testing.T) {
	repos := NewTestScrollerRepos()
	ids := CreateImagesWithOrderedIds(repos.Image, 3)
	r, _ := repos.Scroller.GetAdjacent(ids[2], scroller.NewCriteria(), scroller.ScrollPrevious)
	if r.ImageId != ids[1] {
		t.Fatalf("expected to retrieve next image with id %v got %v", ids[1], r.ImageId)
	}
}

func TestScrollWithCollectionCriteria(t *testing.T) {
	repos := NewTestScrollerRepos()
	firstImage := CreateImageInCollection(repos.Image, repos.Collection,
		im.NewImageId(), "first-collection")
	CreateImageInCollection(repos.Image, repos.Collection,
		im.NewImageId(), "second-collection")

	_, err := repos.Scroller.GetAdjacent(firstImage.Id,
		scroller.NewCriteria(scroller.WithCollection("first-collection")),
		scroller.ScrollPrevious)

	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found err got  %v", err)
	}
}

func TestGettingNextImageInCollection(t *testing.T) {
	repos := NewTestScrollerRepos()
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	ids := CreateImagesWithOrderedIds(repos.Image, 2)
	repos.Collection.Create(*collection)
	repos.Image.AddToCollection(ids[0], collection.Id)
	repos.Image.AddToCollection(ids[1], collection.Id)

	r, _ := repos.Scroller.GetAdjacent(ids[0],
		scroller.NewCriteria(scroller.WithCollection(collection.Name)),
		scroller.ScrollNext)
	if r.Collection != collection.Name {
		t.Fatalf("expected to retrieve next image in collection %v got %v",
			collection.Name, r.Collection)
	}
}
