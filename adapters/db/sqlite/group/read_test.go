package group

import (
	"testing"

	g "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	_, err := repo.Find("non-existing-group")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	repo.Db.Close()
	_, err := repo.Find("a-group")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRetrieve(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	grp := g.NewGroup(g.NewGroupId(), "a-group",
		g.WithDescription("a-description"))
	repo.Create(grp)
	r, err := repo.Find("a-group")
	assert.NoError(t, err, "expected no error on find")
	assert.Equal(t, grp.Name, r.Name)
	assert.Equal(t, grp.Description, r.Description)
	assert.Equal(t, grp.Id, r.Id)

}
