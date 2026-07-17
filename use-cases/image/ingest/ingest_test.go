package ingest

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	ig "github.com/lejeunel/go-image-annotator/modules/ingester"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := grp.NewGroup(grp.NewGroupId(), "a-group")
	itr := NewTestingInteractor(&fk.CollectionRepo{
		Return: clc.NewCollection(clc.NewCollectionId(),
			"a-collection",
			clc.WithGroup(group))},
		WithAuth(&fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), ig.Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor(&fk.CollectionRepo{ErrOnFind: e.ErrNotFound})
	itr.Execute(t.Context(), ig.Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnCollectionExistsCheck(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor(&fk.CollectionRepo{ErrOnExists: e.ErrInternal})
	itr.Execute(t.Context(), ig.Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
