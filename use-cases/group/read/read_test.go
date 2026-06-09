package read

import (
	"github.com/stretchr/testify/assert"
	"testing"

	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestRead(t *testing.T) {
	group := grp.NewGroup(grp.NewGroupId(), "my-group")
	repo := &FakeRepo{Group: group}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	itr.Execute(t.Context(), Request{Name: group.Name}, p)
	assert.Equal(t, Response{Name: group.Name}, p.Got)
}

func TestReadNonExistingShouldFail(t *testing.T) {
	group := grp.NewGroup(grp.NewGroupId(), "my-group")
	repo := &FakeRepo{Group: group}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	req := Request{Name: "non-existing-group"}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}
