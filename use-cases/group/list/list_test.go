package list

import (
	"testing"

	grp "github.com/lejeunel/go-image-annotator/entities/group"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.GroupRepo{ErrOnList: e.ErrInternal})
	itr.Execute(t.Context(), p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestList(t *testing.T) {
	g0 := grp.NewGroup(grp.NewGroupId(), "first-group")
	g1 := grp.NewGroup(grp.NewGroupId(), "second-group")
	repo := &fk.GroupRepo{ReturnList: []grp.Group{g0, g1}}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), p)
	assert.Equal(t, 2, len(p.Got))
}
