package group

import (
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	_, err := CreateGroup(repo, "a-group")
	assert.ErrorIs(t, err, e.ErrInternal, "expected internal error")
}

func TestCreate(t *testing.T) {
	_, err := CreateGroup(NewTestSQLiteGroupRepo(), "a-group")
	assert.NoError(t, err, "expected no error on create but got")

}
