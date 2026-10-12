package profile

import (
	"testing"

	lbr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	err := repo.Delete("a-profile")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteProfile(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	label, _ := lbr.CreateLabel(lbr.NewLabelRepo(db), "a-label")
	profile, _ := CreateProfile(repo, "a-profile", nil, &[]string{label.Name})
	err := repo.Delete(profile.Name)
	assert.NoError(t, err)
}
