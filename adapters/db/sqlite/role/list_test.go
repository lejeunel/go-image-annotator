package role

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnListShouldFail(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	repo.Db.Close()
	_, err := repo.List()
	assert.ErrorIs(t, err, e.ErrInternal)
}
func TestListEmpty(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	_, err := repo.List()
	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	CreateRole(repo, "a-role")
	CreateRole(repo, "another-role")
	cs, err := repo.List()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cs))
}
