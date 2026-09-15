package profile

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreatedProfileExists(t *testing.T) {
	repo := NewProfileRepo(s.NewInMemory())
	profile, _ := CreateProfile(repo, "a-profile")
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

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	err := repo.Delete("a-profile")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteProfile(t *testing.T) {
	repo := NewProfileRepo(s.NewInMemory())
	profile, _ := CreateProfile(repo, "a-profile")
	err := repo.Delete(profile.Name)
	assert.NoError(t, err)
}
