package role

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnUpdateShouldFail(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	repo.Db.Close()
	err := repo.Update(rl.UpdatableModel{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdate(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	role, _ := CreateRole(repo, "a-role")
	newName := "new-role-name"
	newDesc := "new-description"
	err := repo.Update(rl.UpdatableModel{Name: role.Name, NewName: newName, NewDescription: newDesc})
	assert.NoError(t, err)
	r, err := repo.Find(newName)
	assert.NoError(t, err)
	assert.Equal(t, newName, r.Name)
	assert.Equal(t, newDesc, r.Description)
}
