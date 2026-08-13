package collection

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	grr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewCollectionRepo(db)
	db.Close()
	_, err := CreateCollection(repo, "a-collection")
	assert.ErrorIs(t, err, e.ErrInternal, "expected internal error")
}

func TestCreate(t *testing.T) {
	_, err := CreateCollection(NewCollectionRepo(s.NewInMemory()), "a-collection")
	assert.NoError(t, err, "expected no error on create but got")
}

func TestCreateCollectionInGroup(t *testing.T) {
	db := s.NewInMemory()
	groupRepo := grr.NewGroupRepo(db)
	collectionRepo := NewCollectionRepo(db)
	group := grp.NewGroup(grp.NewGroupId(), "a-group")
	groupRepo.Create(group)
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithGroup(group.Name))
	collectionRepo.Create(c)
	r, err := collectionRepo.GetGroup(c.Name)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "a-group", *r)
}

func TestCollectionWithoutGroupFailsWithNotFoundErr(t *testing.T) {
	db := s.NewInMemory()
	clcRepo := NewCollectionRepo(db)
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	clcRepo.Create(collection)
	group, _ := clcRepo.GetGroup("a-collection")
	assert.Nil(t, group)
}
