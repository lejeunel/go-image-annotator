package collection

import (
	"testing"

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
