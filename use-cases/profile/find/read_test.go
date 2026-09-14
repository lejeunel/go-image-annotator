package find

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestReadProfile(t *testing.T) {
	profile := pr.NewProfile(pr.NewProfileId(),
		"my-profile",
		pr.WithDescription("a-description"))
	repo := &fk.ProfileRepo{Return: profile}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), profile.Name, p)
	assert.Equal(t, profile, p.Got)
}

func TestErrorOnFInd(t *testing.T) {
	repo := &fk.ProfileRepo{ErrOnFind: e.ErrNotFound}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), "non-existing-profile", p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}
