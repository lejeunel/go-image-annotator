package group

import (
	"errors"
	"testing"

	g "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	_, err := repo.Find("non-existing-group")
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	repo.Db.Close()
	_, err := repo.Find("a-group")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
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
