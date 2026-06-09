package user

import (
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

var userId = "user@example.com"

func TestRetrieveUserWithNoRole(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, userId)
	r, _ := repo.Find(userId)
	assert.Equal(t, 0, len(r.Roles))
}

func TestCreateUserWithOneRole(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	user := u.NewUser(userId, u.WithRoles([]string{"a-role"}))
	repo.Create(user)
	r, _ := repo.Find(userId)
	assert.Equal(t, 1, len(r.Roles))
}

func TestAssignRoleToExistingUser(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, userId)
	err := repo.AssignRole(userId, "a-role")
	assert.NoError(t, err)
	r, err := repo.Find(userId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.Roles))
}

func TestAssignSameRoleTwice(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, userId)
	repo.AssignRole(userId, "a-role")
	repo.AssignRole(userId, "a-role")
	r, err := repo.Find(userId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.Roles))
}

func TestUnAssignNonExistingRoleShouldFail(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, userId)
	err := repo.UnAssignRole(userId, "non-existing-role")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestUnAssignRole(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, userId)
	repo.AssignRole(userId, "a-role")
	repo.UnAssignRole(userId, "a-role")
	r, _ := repo.Find(userId)
	assert.Equal(t, 0, len(r.Roles))
}
