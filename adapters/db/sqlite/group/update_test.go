package group

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	u "github.com/lejeunel/go-image-annotator/use-cases/group/update"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnUpdateShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	err := repo.Update(u.Model{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdate(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	group, _ := CreateGroup(repo, "a-group")
	newName := "new-group-name"
	newDesc := "new-description"
	err := repo.Update(u.Model{Name: group.Name, NewName: newName, NewDescription: newDesc})
	assert.NoError(t, err)
	r, err := repo.Find(newName)
	assert.NoError(t, err)
	assert.Equal(t, newName, r.Name)
	assert.Equal(t, newDesc, r.Description)
}
