package find

import (
	"github.com/stretchr/testify/assert"
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	fk "github.com/lejeunel/go-image-annotator/use-cases/fakes"
)

func TestReadCollection(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(),
		"my-collection",
		clc.WithDescription("a-description"))
	repo := &fk.CollectionRepo{Return: collection}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), collection.Name, p)
	assert.Equal(t, collection, p.Got)
}

func TestErrorOnFInd(t *testing.T) {
	repo := &fk.CollectionRepo{ErrOnFind: e.ErrNotFound}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), "non-existing-collection", p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}
