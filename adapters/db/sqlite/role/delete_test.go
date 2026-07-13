package role

import (
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	usrRepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatedExists(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	role, _ := CreateRole(repo, "a-role")
	exists, _ := repo.Exists(role.Name)
	assert.True(t, *exists)
}

func TestNonExistingDoesNotExists(t *testing.T) {
	exists, _ := NewTestSQLiteRoleRepo().Exists("non-existing-role")
	assert.False(t, *exists)
}

func TestInternalErrOnExistsShouldFail(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	repo.Db.Close()
	_, err := repo.Exists("")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	repo.Db.Close()
	err := repo.Delete("a-role")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDelete(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	grp, _ := CreateRole(repo, "a-role")
	err := repo.Delete(grp.Name)
	assert.Nil(t, err)
}

func TestRoleUsedByUser(t *testing.T) {
	db := s.NewSQLiteDB(":memory:")
	usrRepo := usrRepo.NewSQLiteUserRepo(db)
	roleRepo := NewSQLiteRoleRepo(db)
	role, _ := CreateRole(roleRepo, "a-role")
	user := usr.NewUser("user@mail.com", usr.WithRoles([]string{"a-role"}))
	usrRepo.Create(user)
	isPopulated, err := roleRepo.IsAssigned(role.Name)
	assert.NoError(t, err)
	assert.True(t, *isPopulated)
}
