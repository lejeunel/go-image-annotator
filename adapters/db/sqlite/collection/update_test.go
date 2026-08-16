package collection

import (
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	gr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCollectionUpdateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewCollectionRepo(db)
	db.Close()
	err := repo.Update(clc.UpdateModel{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func Setup() (CollectionRepo, clc.Collection, gr.GroupRepo, g.Group) {
	db := s.NewInMemory()
	clcRepo := NewCollectionRepo(db)
	grpRepo := gr.NewGroupRepo(db)
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection",
		clc.WithDescription("a-description"), clc.WithCreatedAt(time.Now()))
	clcRepo.Create(collection)

	group := g.NewGroup(g.NewGroupId(), "my-group")
	grpRepo.Create(group)
	return clcRepo, collection, grpRepo, group
}

func TestUpdateNameAndDescription(t *testing.T) {
	clcRepo, collection, _, _ := Setup()
	req := clc.UpdateModel{
		Name: collection.Name, NewName: "new-collection-name",
		NewDescription: "new-description",
	}
	err := clcRepo.Update(req)
	assert.NoError(t, err)
	r, err := clcRepo.Find(req.NewName)
	assert.NoError(t, err)
	assert.Equal(t, req.NewName, r.Name)
	assert.Equal(t, req.NewDescription, r.Description)
}

func TestUpdateGroupFromPublic(t *testing.T) {
	clcRepo, collection, _, group := Setup()
	req := clc.UpdateModel{
		Name: collection.Name, NewName: collection.Name,
		NewDescription: collection.Description,
		NewGroup:       &group.Name,
	}
	err := clcRepo.Update(req)
	assert.NoError(t, err)
	r, err := clcRepo.Find(req.NewName)
	assert.NoError(t, err)
	assert.NotNil(t, r.Group)
}

func TestUpdateGroupToPublic(t *testing.T) {
	clcRepo, collection, _, group := Setup()
	req := clc.UpdateModel{
		Name: collection.Name, NewName: collection.Name,
		NewDescription: collection.Description,
		NewGroup:       &group.Name,
	}
	clcRepo.Update(req)
	req.NewGroup = nil
	clcRepo.Update(req)
	r, err := clcRepo.Find(req.NewName)
	assert.NoError(t, err)
	assert.Nil(t, r.Group)
}

func TestUpdateAndListCollections(t *testing.T) {
	clcRepo, collection, _, group := Setup()
	req := clc.UpdateModel{
		Name: collection.Name, NewName: "new-collection-name",
		NewDescription: "new-description",
		NewGroup:       &group.Name,
	}
	err := clcRepo.Update(req)
	assert.NoError(t, err)
	r, err := clcRepo.List(pagination.PaginationParams{Page: 1, PageSize: 1})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r))
	assert.NotNil(t, r[0].Group)
}
