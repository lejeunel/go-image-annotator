package find

import (
	"github.com/stretchr/testify/assert"
	"testing"

	grp "github.com/lejeunel/go-image-annotator/entities/group"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestRead(t *testing.T) {
	group := grp.NewGroup(grp.NewGroupId(), "my-group")
	repo := &fk.GroupRepo{Return: group}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), group.Name, p)
	assert.Equal(t, group, p.Got)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.GroupRepo{ErrOnFind: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
}
