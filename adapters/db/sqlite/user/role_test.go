package user

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	roleRepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
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
	db := s.NewSQLiteDB(":memory:")
	repo := NewSQLiteUserRepo(db)
	roleRepo := roleRepo.NewSQLiteRoleRepo(db)
	user := u.NewUser(userId, u.WithRoles([]string{"a-role"}))
	roleRepo.Create(r.NewRole(r.NewRoleId(), "a-role"))
	repo.Create(user)
	r, _ := repo.Find(userId)
	assert.Equal(t, 1, len(r.Roles))
}

func TestAssignRoleToExistingUser(t *testing.T) {
	db := s.NewSQLiteDB(":memory:")
	repo := NewSQLiteUserRepo(db)
	roleRepo := roleRepo.NewSQLiteRoleRepo(db)
	CreateUser(repo, userId)
	roleRepo.Create(r.NewRole(r.NewRoleId(), "a-role"))
	err := repo.AssignRole(userId, "a-role")
	assert.NoError(t, err)
	r, err := repo.Find(userId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.Roles))
}

func TestAssignSameRoleTwice(t *testing.T) {
	db := s.NewSQLiteDB(":memory:")
	repo := NewSQLiteUserRepo(db)
	CreateUser(repo, userId)
	roleRepo := roleRepo.NewSQLiteRoleRepo(db)
	roleRepo.Create(r.NewRole(r.NewRoleId(), "a-role"))
	err := repo.AssignRole(userId, "a-role")
	err = repo.AssignRole(userId, "a-role")
	r, err := repo.Find(userId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.Roles))
}

func TestUnAssignRole(t *testing.T) {
	db := s.NewSQLiteDB(":memory:")
	repo := NewSQLiteUserRepo(db)
	CreateUser(repo, userId)
	roleRepo := roleRepo.NewSQLiteRoleRepo(db)
	roleRepo.Create(r.NewRole(r.NewRoleId(), "a-role"))
	repo.AssignRole(userId, "a-role")
	repo.UnAssignRole(userId, "a-role")
	r, _ := repo.Find(userId)
	assert.Equal(t, 0, len(r.Roles))
}
