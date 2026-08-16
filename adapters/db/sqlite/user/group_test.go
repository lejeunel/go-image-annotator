package user

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	grpr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveUserWithNoGroup(t *testing.T) {
	repo := NewUserRepo(s.NewInMemory())
	CreateUser(repo, userId)
	r, _ := repo.Find(userId)
	assert.Equal(t, 0, len(r.Groups))
}

func TestCreateUserWithOneGroup(t *testing.T) {
	group := g.NewGroup(g.NewGroupId(), "a-group")
	db := s.NewInMemory()
	userRepo := NewUserRepo(db)
	groupRepo := grpr.NewGroupRepo(db)
	groupRepo.Create(group)
	user := u.NewUser(userId, u.WithGroups([]string{group.Name}))
	err := userRepo.Create(user)
	assert.NoError(t, err)
	r, err := userRepo.Find(user.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.Groups))
}

func TestAssignToNewGroup(t *testing.T) {
	db := s.NewInMemory()
	userRepo := NewUserRepo(db)
	groupRepo := grpr.NewGroupRepo(db)

	group := g.NewGroup(g.NewGroupId(), "a-group")
	groupRepo.Create(group)
	user := u.NewUser(userId, u.WithGroups([]string{group.Name}))
	userRepo.Create(user)

	g0 := g.NewGroup(g.NewGroupId(), "a-new-group")
	groupRepo.Create(g0)
	g1 := g.NewGroup(g.NewGroupId(), "another-new-group")
	groupRepo.Create(g1)

	err := userRepo.SetGroups(user.Id, []string{g0.Name, g1.Name})
	assert.NoError(t, err)
	r, _ := userRepo.Find(user.Id)
	assert.Equal(t, 2, len(r.Groups))
}
