package user

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	grpr "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/group"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	// e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveUserWithNoGroup(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, userId)
	r, _ := repo.Find(userId)
	assert.Equal(t, 0, len(r.Groups))
}

func TestCreateUserWithOneGroup(t *testing.T) {
	group := g.NewGroup(g.NewGroupId(), "a-group")
	db := s.NewSQLiteDB(":memory:")
	userRepo := NewSQLiteUserRepo(db)
	groupRepo := grpr.NewSQLiteGroupRepo(db)
	groupRepo.Create(group)
	user := u.NewUser(userId, u.WithGroups([]string{group.Name}))
	err := userRepo.Create(user)
	assert.NoError(t, err)
	r, err := userRepo.Find(user.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.Groups))
}

func TestAssignSameGroupTwice(t *testing.T) {
	group := g.NewGroup(g.NewGroupId(), "a-group")
	db := s.NewSQLiteDB(":memory:")
	userRepo := NewSQLiteUserRepo(db)
	groupRepo := grpr.NewSQLiteGroupRepo(db)
	groupRepo.Create(group)
	user := u.NewUser(userId)
	userRepo.Create(user)
	userRepo.AssignToGroup(user.Id, group.Name)
	userRepo.AssignToGroup(user.Id, group.Name)
	r, _ := userRepo.Find(user.Id)
	assert.Equal(t, 1, len(r.Groups))
}

func TestUnAssignGroup(t *testing.T) {
	group := g.NewGroup(g.NewGroupId(), "a-group")
	db := s.NewSQLiteDB(":memory:")
	userRepo := NewSQLiteUserRepo(db)
	groupRepo := grpr.NewSQLiteGroupRepo(db)
	groupRepo.Create(group)
	user := u.NewUser(userId, u.WithGroups([]string{group.Name}))
	userRepo.Create(user)
	userRepo.UnAssignFromGroup(user.Id, group.Name)
	r, _ := userRepo.Find(user.Id)
	assert.Equal(t, 0, len(r.Groups))
}
