package list

import (
	"testing"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnList: e.ErrInternal})
	itr.Execute(t.Context(), p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestListProfiles(t *testing.T) {
	profile := pr.NewProfile(pr.NewProfileId(), "my-profile")
	repo := &fk.ProfileRepo{Return: profile}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), p)

	assert.Equal(t, 1, len(p.Got))
	assert.Equal(t, profile.Name, p.Got[0])
}
