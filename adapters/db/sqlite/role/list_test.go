package role

import (
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnListShouldFail(t *testing.T) {
	repo := NewTestRoleRepo()
	repo.Db.Close()
	_, err := repo.List()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestListEmpty(t *testing.T) {
	repo := NewTestRoleRepo()
	_, err := repo.List()
	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	repo := NewTestRoleRepo()
	CreateRole(repo, "a-role")
	CreateRole(repo, "another-role")
	cs, err := repo.List()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cs))
}
