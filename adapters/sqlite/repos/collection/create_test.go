package collection

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	grr "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/group"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	_, err := CreateCollection(repo, "a-collection")
	assert.ErrorIs(t, err, e.ErrInternal, "expected internal error")
}

func TestCreate(t *testing.T) {
	_, err := CreateCollection(NewTestSQLiteCollectionRepo(), "a-collection")
	assert.NoError(t, err, "expected no error on create but got")
}

func TestCreateCollectionInGroup(t *testing.T) {
	db := s.NewSQLiteDB(":memory:")
	groupRepo := grr.NewSQLiteGroupRepo(db)
	collectionRepo := NewSQLiteCollectionRepo(db)
	group := grp.NewGroup(grp.NewGroupId(), "a-group")
	groupRepo.Create(group)
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithGroup(group))
	collectionRepo.Create(c)
	r, err := groupRepo.GroupOfCollection(c.Name)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "a-group", *r)
}
