package profile

import (
	"testing"

	gr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnProfileUpdateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	err := repo.Update(pr.UpdateModel{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func Setup() (ProfileRepo, pr.Profile, gr.GroupRepo, g.Group) {
	db := s.NewInMemory()
	prRepo := NewProfileRepo(db)
	grpRepo := gr.NewGroupRepo(db)
	profile := pr.NewProfile(pr.NewProfileId(), "my-profile",
		pr.WithDescription("a-description"))
	prRepo.Create(profile)

	group := g.NewGroup(g.NewGroupId(), "my-group")
	grpRepo.Create(group)
	return prRepo, profile, grpRepo, group
}

func TestUpdateNameAndDescription(t *testing.T) {
	prRepo, profile, _, _ := Setup()
	req := pr.UpdateModel{
		Name: profile.Name, NewName: "new-profile-name",
		NewDescription: "new-description",
	}
	err := prRepo.Update(req)
	assert.NoError(t, err)
	r, err := prRepo.Find(req.NewName)
	assert.NoError(t, err)
	assert.Equal(t, req.NewName, r.Name)
	assert.Equal(t, req.NewDescription, r.Description)
}

func TestUpdateGroupFromPublic(t *testing.T) {
	prRepo, profile, _, group := Setup()
	req := pr.UpdateModel{
		Name: profile.Name, NewName: profile.Name,
		NewDescription: profile.Description,
		NewGroup:       &group.Name,
	}
	err := prRepo.Update(req)
	assert.NoError(t, err)
	r, err := prRepo.Find(req.NewName)
	assert.NoError(t, err)
	assert.NotNil(t, r.Group)
}

func TestUpdateGroupToPublic(t *testing.T) {
	prRepo, profile, _, group := Setup()
	req := pr.UpdateModel{
		Name: profile.Name, NewName: profile.Name,
		NewDescription: profile.Description,
		NewGroup:       &group.Name,
	}
	prRepo.Update(req)
	req.NewGroup = nil
	prRepo.Update(req)
	r, err := prRepo.Find(req.NewName)
	assert.NoError(t, err)
	assert.Nil(t, r.Group)
}

func TestUpdateAndListProfiles(t *testing.T) {
	prRepo, profile, _, group := Setup()
	req := pr.UpdateModel{
		Name: profile.Name, NewName: "new-profile-name",
		NewDescription: "new-description",
		NewGroup:       &group.Name,
	}
	err := prRepo.Update(req)
	assert.NoError(t, err)
	r, err := prRepo.List(pagination.PaginationParams{Page: 1, PageSize: 1})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r))
	assert.NotNil(t, r[0].Group)
}
