package user

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	grr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	rlr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func CreateUser(repo SQLiteUserRepo, id string, opts ...u.Option) (*u.User, error) {
	user := u.NewUser(id, opts...)
	if err := repo.Create(user); err != nil {
		return nil, err
	}
	return &user, nil

}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	repo.Db.Close()
	_, err := CreateUser(repo, userId)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCreateAddsCount(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	_, err := CreateUser(repo, userId)
	assert.NoError(t, err)
	count, err := repo.Count()
	assert.Equal(t, int64(1), count)
	assert.NoError(t, err)
}

func TestNoCreatedUserDoNotExist(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	exists, err := repo.Exists(userId)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestPersonalAccessTokenHash(t *testing.T) {
	hash := []byte("pat-hash")
	repo := NewSQLiteUserRepo(s.NewInMemory())
	user := u.NewUser("user@example.com",
		u.WithHashedPersonalAccessToken(hash))
	repo.Create(user)
	r, err := repo.Find("user@example.com")
	assert.NoError(t, err)
	assert.Equal(t, hash, r.HashPAT)
}

func TestPasswordHash(t *testing.T) {
	hash := []byte("password-hash")
	repo := NewSQLiteUserRepo(s.NewInMemory())
	user := u.NewUser("user@example.com",
		u.WithHashedPassword(hash))
	repo.Create(user)
	r, err := repo.Find("user@example.com")
	assert.NoError(t, err)
	assert.Equal(t, hash, r.HashPassword)
}

func TestCreateAdminInGroup(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteUserRepo(db)

	roleRepo := rlr.NewSQLiteRoleRepo(db)
	role := r.NewRole(r.NewRoleId(), "admin")
	roleRepo.Create(role)

	groupRepo := grr.NewSQLiteGroupRepo(db)
	group := g.NewGroup(g.NewGroupId(), "my-group")
	groupRepo.Create(group)

	user := u.NewUser("user@example.com",
		u.WithRoles([]string{"admin"}),
		u.WithGroups([]string{"my-group"}))
	err := repo.Create(user)
	assert.NoError(t, err)
	r, err := repo.Find("user@example.com")
	assert.NoError(t, err)
	assert.Contains(t, r.Roles, "admin")
	assert.Contains(t, r.Groups, "my-group")

}
