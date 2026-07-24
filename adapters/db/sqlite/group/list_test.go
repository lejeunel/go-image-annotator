package group

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnCountShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	_, err := repo.Count()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCount(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	count, _ := repo.Count()
	assert.Equal(t, 1, int(*count))
}

func TestInternalErrOnListShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	_, err := repo.List()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestList(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	CreateGroup(repo, "another-group")
	cs, err := repo.List()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cs))
	assert.False(t, cs[0].Name == cs[1].Name)
}
