package label

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatedLabelExists(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	label, _ := CreateLabel(repo, "a-label")
	exists, _ := repo.Exists(label.Name)
	assert.True(t, exists)
}

func TestNonExistingLabelDoesNotExists(t *testing.T) {
	exists, _ := NewTestSQLiteLabelRepo().Exists("non-existing-label")
	assert.False(t, exists)
}

func TestInternalErrOnLabelExistsShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := repo.Exists("")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	err := repo.Delete("a-label")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteLabel(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	label, _ := CreateLabel(repo, "a-label")
	err := repo.Delete(label.Name)
	assert.NoError(t, err)
}
