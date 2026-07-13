package role

import (
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteRoleRepo()
	repo.Db.Close()
	_, err := CreateRole(repo, "a-role")
	assert.ErrorIs(t, err, e.ErrInternal, "expected internal error")
}

func TestCreate(t *testing.T) {
	_, err := CreateRole(NewTestSQLiteRoleRepo(), "a-role")
	assert.NoError(t, err, "expected no error on create but got")

}
