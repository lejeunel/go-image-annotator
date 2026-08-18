package scroll

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func Setup() Interactor {
	return New(&fk.ImageRepo{}, &fk.FilterValidator{}, &fk.FilterValidator{})

}

func TestErrOnInvalidImageId(t *testing.T) {
	p := &FakePresenter{}
	itr := Setup()
	itr.Execute(t.Context(), Request{CurrentImageId: "invalid-image-id"}, p)
	assert.Error(t, p.GotErr)
}

func TestErrOnCurrentImageExistsShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := Setup()
	itr.ImageRepo = &fk.ImageRepo{ErrOnImageExists: e.ErrInternal}
	itr.Execute(t.Context(), Request{CurrentImageId: im.NewImageId().String()}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}
func TestNonExistingCurrentImageExistsShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := Setup()
	itr.ImageRepo = &fk.ImageRepo{ImageMissing: true}
	itr.Execute(t.Context(), Request{CurrentImageId: im.NewImageId().String()}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
}
func TestErrOnFilteringStrValidation(t *testing.T) {
	p := &FakePresenter{}
	itr := Setup()
	v := fk.FilterValidator{Err: e.ErrValidation}
	itr.FilterValidator = &v
	queryStr := "collection:\"my-collection\""
	itr.Execute(t.Context(), Request{CurrentImageId: im.NewImageId().String(),
		FilterStr: queryStr}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
	assert.Equal(t, v.Got, queryStr)
}

func TestErrOnOrderingStrValidation(t *testing.T) {
	p := &FakePresenter{}
	itr := Setup()
	v := fk.FilterValidator{Err: e.ErrValidation}
	itr.OrderingValidator = &v
	orderStr := "ingested_at:desc"
	itr.Execute(t.Context(), Request{CurrentImageId: im.NewImageId().String(),
		OrderStr: orderStr}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
	assert.Equal(t, v.Got, orderStr)
}

func TestErrOnGetAdjacent(t *testing.T) {
	p := &FakePresenter{}
	itr := Setup()
	itr.ImageRepo = &fk.ImageRepo{ErrOnGetAdjacent: e.ErrInternal}
	itr.Execute(t.Context(), Request{CurrentImageId: im.NewImageId().String()}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestGetAdjacent(t *testing.T) {
	p := &FakePresenter{}
	itr := Setup()
	collection := "my-collection"
	nextId, prevId := im.NewImageId(), im.NewImageId()
	next := im.BaseImage{ImageId: nextId, Collection: collection}
	prev := im.BaseImage{ImageId: prevId, Collection: collection}
	adj := im.AdjacentImages{Next: &next, Prev: &prev}
	itr.ImageRepo = &fk.ImageRepo{Adjacent: adj}
	itr.Execute(t.Context(), Request{CurrentImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, nextId.String(), p.Got.Adj.Next.ImageId.String())
	assert.Equal(t, prevId.String(), p.Got.Adj.Prev.ImageId.String())
}
