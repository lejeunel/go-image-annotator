package label

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	l "github.com/lejeunel/go-image-annotator/use-cases/label/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnLabelCountShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := repo.Count()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCountLabels(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	CreateLabel(repo, "a-label")
	count, _ := repo.Count()
	assert.Equal(t, 1, int(count))
}

func TestInternalErrOnLabelListShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := repo.List(l.Request{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestListLabels(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	CreateLabel(repo, "a-label")
	CreateLabel(repo, "another-label")
	labels, err := repo.List(l.Request{Page: 1, PageSize: 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(labels))
	assert.NotEqual(t, labels[0].Name, labels[1].Name)
}
