package profile

import (
	"testing"

	gr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	_, err := CreateProfile(repo, "a-profile", nil)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCreatedProfileExists(t *testing.T) {
	repo := NewProfileRepo(s.NewInMemory())
	profile, _ := CreateProfile(repo, "a-profile", nil)
	exists, _ := repo.Exists(profile.Name)
	assert.True(t, exists)
}

func TestNonExistingProfileDoesNotExists(t *testing.T) {
	exists, _ := NewProfileRepo(s.NewInMemory()).Exists("non-existing-profile")
	assert.False(t, exists)
}

func TestInternalErrOnProfileExistsShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	_, err := repo.Exists("")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestProfileWithoutGroup(t *testing.T) {
	db := s.NewInMemory()
	profileRepo := NewProfileRepo(db)
	pr, err := CreateProfile(profileRepo, "a-profile", nil)
	assert.NoError(t, err)
	group, _ := profileRepo.GetGroup(pr.Name)
	assert.Nil(t, group)
}

func TestCreateProfileInGroup(t *testing.T) {
	db := s.NewInMemory()
	groupRepo := gr.NewGroupRepo(db)
	group := grp.NewGroup(grp.NewGroupId(), "a-group")
	groupRepo.Create(group)

	profileRepo := NewProfileRepo(db)
	p, err := CreateProfile(profileRepo, "a-profile", &group.Name)
	assert.NoError(t, err)
	r, err := profileRepo.GetGroup(p.Name)
	assert.NoError(t, err)
	assert.Equal(t, "a-group", *r)
}
