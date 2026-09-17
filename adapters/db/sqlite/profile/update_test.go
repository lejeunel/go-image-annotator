package profile

import (
	"testing"

	gr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	lr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnProfileUpdateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	err := repo.Update(pr.UpdateModel{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func Setup(labels *[]string) (ProfileRepo, pr.Profile, gr.GroupRepo, lr.LabelRepo, g.Group) {
	db := s.NewInMemory()
	prRepo := NewProfileRepo(db)
	grpRepo := gr.NewGroupRepo(db)
	labelRepo := lr.NewLabelRepo(db)
	profile := pr.NewProfile(pr.NewProfileId(), "my-profile",
		pr.WithDescription("a-description"))

	if labels != nil {
		for _, label := range *labels {
			if err := labelRepo.Create(lbl.NewLabel(lbl.NewLabelId(), label)); err != nil {
				panic(err)
			}
		}
		profile.Labels = *labels
	}

	prRepo.Create(profile)
	group := g.NewGroup(g.NewGroupId(), "my-group")
	grpRepo.Create(group)
	return prRepo, profile, grpRepo, labelRepo, group
}

func TestUpdateNameAndDescription(t *testing.T) {
	prRepo, profile, _, _, _ := Setup(nil)
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
	prRepo, profile, _, _, group := Setup(nil)
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
	prRepo, profile, _, _, group := Setup(nil)
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

func TestUpdateLabels(t *testing.T) {
	originalLabel := "a-label"
	prRepo, profile, _, labelRepo, _ := Setup(&[]string{originalLabel})
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "new-label")
	labelRepo.Create(newLabel)
	req := pr.UpdateModel{
		Name: profile.Name, NewName: profile.Name,
		NewDescription: profile.Description,
		NewGroup:       profile.Group,
		NewLabels:      []string{newLabel.Name},
	}
	err := prRepo.Update(req)
	assert.NoError(t, err)
	upd, _ := prRepo.Find(profile.Name)
	assert.Equal(t, 1, len(upd.Labels))
	assert.Equal(t, newLabel.Name, upd.Labels[0])
}
