package role

import (
	"testing"

	ro "github.com/lejeunel/go-image-annotator/entities/role"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	CreateRole(repo, "a-role")
	_, err := repo.Find("non-existing-role")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	CreateRole(repo, "a-role")
	repo.Db.Close()
	_, err := repo.Find("a-role")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRetrieve(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	role := ro.NewRole(ro.NewRoleId(), "a-role",
		ro.WithDescription("a-description"))
	repo.Create(role)
	r, err := repo.Find("a-role")
	assert.NoError(t, err, "expected no error on find")
	assert.Equal(t, role.Name, r.Name)
	assert.Equal(t, role.Description, r.Description)
	assert.Equal(t, role.Id, r.Id)

}
