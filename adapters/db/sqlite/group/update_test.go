package group

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnUpdateShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	err := repo.Update(grp.UpdateModel{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdate(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	group, _ := CreateGroup(repo, "a-group")
	newName := "new-group-name"
	newDesc := "new-description"
	err := repo.Update(grp.UpdateModel{Name: group.Name, NewName: newName, NewDescription: newDesc})
	assert.NoError(t, err)
	r, err := repo.Find(newName)
	assert.NoError(t, err)
	assert.Equal(t, newName, r.Name)
	assert.Equal(t, newDesc, r.Description)
}
