package group

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	l "github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnCountShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	_, err := repo.Count()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCount(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	count, _ := repo.Count()
	assert.Equal(t, 1, int(*count))
}

func TestInternalErrOnListShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	_, err := repo.List(l.Request{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestList(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	CreateGroup(repo, "another-group")
	cs, err := repo.List(l.Request{Page: 1, PageSize: 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cs))
	assert.False(t, cs[0].Name == cs[1].Name)
}
