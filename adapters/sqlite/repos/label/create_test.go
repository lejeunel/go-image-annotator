package label

import (
	"testing"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func CreateLabel(repo *SQLiteLabelRepo, name string) (*lbl.Label, error) {
	label := lbl.NewLabel(lbl.NewLabelId(), name, lbl.WithDescription("a-description"))
	if err := repo.Create(label); err != nil {
		return nil, err
	}
	return &label, nil

}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := CreateLabel(repo, "a-label")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCreateAddsCount(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	_, err := CreateLabel(repo, "a-label")
	assert.NoError(t, err)
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, int(count))
}
