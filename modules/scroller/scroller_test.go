package scroller

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrOnInvalidImageId(t *testing.T) {
	s := New(&FakeRepo{})
	_, err := s.Init("invalid-image-id")
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestErrOnCurrentImageShouldFail(t *testing.T) {
	s := New(&FakeRepo{ErrOnImageExists: true, Err: e.ErrNotFound})
	_, err := s.Init(im.NewImageId().String())
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	s := New(&FakeRepo{ErrOnCollectionExists: true, Err: e.ErrNotFound})
	_, err := s.Init(im.NewImageId().String(), WithCollection("non-existing-collection"))
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestSingleImageHasNoNextImage(t *testing.T) {
	s := New(&FakeRepo{})
	state, _ := s.Init(im.NewImageId().String())
	assert.Nil(t, state.Next)
}

func TestSingleImageHasNoPreviousImage(t *testing.T) {
	s := New(&FakeRepo{})
	state, _ := s.Init(im.NewImageId().String())
	assert.Nil(t, state.Previous)
}

func TestNextImage(t *testing.T) {
	next := &im.BaseImage{ImageId: im.NewImageId().String()}
	s := New(&FakeRepo{NextImage: next})
	state, _ := s.Init(im.NewImageId().String())
	assert.NotNil(t, state.Next)
	assert.Equal(t, next.ImageId, state.Next.ImageId)
}

func TestPreviousImage(t *testing.T) {
	prev := &im.BaseImage{ImageId: im.NewImageId().String()}
	s := New(&FakeRepo{PreviousImage: prev})
	state, _ := s.Init(im.NewImageId().String())
	assert.NotNil(t, state.Previous)
	assert.Equal(t, prev.ImageId, state.Previous.ImageId)
}
