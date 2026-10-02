package role

import (
	"testing"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMethods(t *testing.T) {
	methods := []string{"a-method"}
	repo := NewTestRoleRepo()
	CreateRole(repo, "a-role")
	req := rl.UpdatableModel{
		Name: "a-role", NewName: "a-role",
		NewMethods: []string{"a-method"},
	}
	err := repo.Update(req)
	assert.NoError(t, err)
	role, err := repo.Find("a-role")
	assert.NoError(t, err)
	assert.Equal(t, methods, role.Methods)
}
