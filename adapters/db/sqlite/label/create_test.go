package label

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewLabelRepo(db)
	db.Close()
	_, err := CreateLabel(repo, "a-label")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCreateAddsCount(t *testing.T) {
	repo := NewLabelRepo(s.NewInMemory())
	_, err := CreateLabel(repo, "a-label")
	assert.NoError(t, err)
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, int(count))
}
