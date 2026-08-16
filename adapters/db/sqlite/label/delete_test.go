package label

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreatedLabelExists(t *testing.T) {
	repo := NewLabelRepo(s.NewInMemory())
	label, _ := CreateLabel(repo, "a-label")
	exists, _ := repo.Exists(label.Name)
	assert.True(t, exists)
}

func TestNonExistingLabelDoesNotExists(t *testing.T) {
	exists, _ := NewLabelRepo(s.NewInMemory()).Exists("non-existing-label")
	assert.False(t, exists)
}

func TestInternalErrOnLabelExistsShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewLabelRepo(db)
	db.Close()
	_, err := repo.Exists("")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewLabelRepo(db)
	db.Close()
	err := repo.Delete("a-label")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteLabel(t *testing.T) {
	repo := NewLabelRepo(s.NewInMemory())
	label, _ := CreateLabel(repo, "a-label")
	err := repo.Delete(label.Name)
	assert.NoError(t, err)
}
