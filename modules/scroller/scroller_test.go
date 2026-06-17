package scroller

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrOnInvalidImageId(t *testing.T) {
	p := &FakePresenter{}
	s := New(&FakeRepo{})
	s.Init("invalid-image-id", p)
	assert.Error(t, p.GotErr)
}

func TestErrOnCurrentImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	s := New(&FakeRepo{ErrOnImageExists: true, Err: e.ErrNotFound})
	s.Init(im.NewImageId().String(), p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	s := New(&FakeRepo{ErrOnCollectionExists: true, Err: e.ErrNotFound})
	s.Init(im.NewImageId().String(), p, WithCollection("non-existing-collection"))
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
}

func TestSingleImageHasNoNextImage(t *testing.T) {
	p := &FakePresenter{}
	s := New(&FakeRepo{})
	s.Init(im.NewImageId().String(), p)
	assert.Nil(t, p.GotState.Next)
}

func TestSingleImageHasNoPreviousImage(t *testing.T) {
	p := &FakePresenter{}
	s := New(&FakeRepo{})
	s.Init(im.NewImageId().String(), p)
	assert.Nil(t, p.GotState.Previous)
}

func TestNextImage(t *testing.T) {
	p := &FakePresenter{}
	next := &im.BaseImage{ImageId: im.NewImageId().String()}
	s := New(&FakeRepo{NextImage: next})
	s.Init(im.NewImageId().String(), p)
	assert.NotNil(t, p.GotState.Next)
	assert.Equal(t, next.ImageId, p.GotState.Next.ImageId)
}

func TestPreviousImage(t *testing.T) {
	p := &FakePresenter{}
	prev := &im.BaseImage{ImageId: im.NewImageId().String()}
	s := New(&FakeRepo{PreviousImage: prev})
	s.Init(im.NewImageId().String(), p)
	assert.NotNil(t, p.GotState.Previous)
	assert.Equal(t, prev.ImageId, p.GotState.Previous.ImageId)
}
