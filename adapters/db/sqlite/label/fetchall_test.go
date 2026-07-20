package label

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnFetchAll(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteLabelRepo(db)
	db.Close()
	_, err := repo.FetchAll()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestFetchAll(t *testing.T) {
	repo := NewSQLiteLabelRepo(s.NewInMemory())
	repo.Create(l.NewLabel(l.NewLabelId(), "first-label"))
	repo.Create(l.NewLabel(l.NewLabelId(), "second-label"))
	labels, _ := repo.FetchAll()
	assert.Equal(t, 2, len(labels))
}
