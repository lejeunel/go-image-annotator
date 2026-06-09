package label

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnFetchAll(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := repo.FetchAll()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestFetchAll(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Create(l.NewLabel(l.NewLabelId(), "first-label"))
	repo.Create(l.NewLabel(l.NewLabelId(), "second-label"))
	labels, _ := repo.FetchAll()
	assert.Equal(t, 2, len(labels))
}
