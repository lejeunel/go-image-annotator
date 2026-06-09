package label

import (
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	CreateLabel(repo, "a-label")
	_, err := repo.FindLabel("non-existing-label")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	CreateLabel(repo, "a-label")
	repo.Db.Close()
	_, err := repo.FindLabel("a-label")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRetrieve(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	label, _ := CreateLabel(repo, "a-label")
	r, err := repo.FindLabel("a-label")
	assert.NoError(t, err)
	assert.Equal(t, r.Name, label.Name)
	assert.Equal(t, r.Description, label.Description)
	assert.Equal(t, r.Id, label.Id)
}
