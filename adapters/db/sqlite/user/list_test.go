package user

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	rlr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	"github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCountShouldFail(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	repo.Db.Close()
	_, err := repo.Count()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCount(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	CreateUser(repo, "user@example.com")
	count, _ := repo.Count()
	assert.Equal(t, 1, int(count))
}

func TestInternalErrOnListShouldFail(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	repo.Db.Close()
	_, err := repo.List(pag.PaginationParams{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestList(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteUserRepo(db)
	roleRepo := rlr.NewSQLiteRoleRepo(db)
	roleRepo.Create(r.NewRole(r.NewRoleId(), "admin"))
	CreateUser(repo, "user@example.com")
	_, err := CreateUser(repo, "another-user@example.com", user.WithRoles([]string{"admin"}))
	users, err := repo.List(pag.PaginationParams{Page: 1, PageSize: 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
	assert.False(t, users[0].Id == users[1].Id)
	assert.True(t, users[1].IsAdmin())
}
