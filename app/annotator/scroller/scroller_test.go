package scroller

import (
	"errors"
	"testing"

	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

func TestErrOnCurrentImageShouldFail(t *testing.T) {
	s := New(&FakeRepo{ErrOnImageExists: true, Err: e.ErrNotFound})
	_, err := s.Init(im.NewImageId())
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	s := New(&FakeRepo{ErrOnCollectionExists: true, Err: e.ErrNotFound})
	_, err := s.Init(im.NewImageId(), WithCollection("non-existing-collection"))
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestSingleImageHasNoNextImage(t *testing.T) {
	s := New(&FakeRepo{})
	state, _ := s.Init(im.NewImageId())
	if state.Next != nil {
		t.Fatal("expected no next image")
	}
}

func TestSingleImageHasNoPreviousImage(t *testing.T) {
	s := New(&FakeRepo{})
	state, _ := s.Init(im.NewImageId())
	if state.Previous != nil {
		t.Fatal("expected no next image")
	}
}

func TestNextImage(t *testing.T) {
	next := &im.BaseImage{ImageId: im.NewImageId()}
	s := New(&FakeRepo{NextImage: next})
	state, _ := s.Init(im.NewImageId())
	if state.Next == nil {
		t.Fatal("expected to get one next image")
	}
	if state.Next.ImageId != next.ImageId {
		t.Fatalf("expected to get next image with id %v, got %v", next.ImageId, state.Next.ImageId)
	}
}

func TestPreviousImage(t *testing.T) {
	prev := &im.BaseImage{ImageId: im.NewImageId()}
	s := New(&FakeRepo{PreviousImage: prev})
	state, _ := s.Init(im.NewImageId())
	if state.Previous == nil {
		t.Fatal("expected to get one previous image")
	}
	if state.Previous.ImageId != prev.ImageId {
		t.Fatalf("expected to get previous image with id %v, got %v", prev.ImageId, state.Previous.ImageId)
	}
}
