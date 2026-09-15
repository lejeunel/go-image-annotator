package profile

import (
	"testing"

	grr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	_, err := CreateProfile(repo, "a-profile")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCreate(t *testing.T) {
	_, err := CreateProfile(NewProfileRepo(s.NewInMemory()), "a-profile")
	assert.NoError(t, err)
}

func TestCreateProfileInGroup(t *testing.T) {
	db := s.NewInMemory()
	groupRepo := grr.NewGroupRepo(db)
	profileRepo := NewProfileRepo(db)
	group := grp.NewGroup(grp.NewGroupId(), "a-group")
	groupRepo.Create(group)
	c := pr.NewProfile(pr.NewProfileId(), "a-profile",
		pr.WithGroup(group.Name))
	profileRepo.Create(c)
	r, err := profileRepo.GetGroup(c.Name)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "a-group", *r)
}

func TestProfileWithoutGroupFailsWithNotFoundErr(t *testing.T) {
	db := s.NewInMemory()
	profileRepo := NewProfileRepo(db)
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	profileRepo.Create(profile)
	group, _ := profileRepo.GetGroup("a-profile")
	assert.Nil(t, group)
}
