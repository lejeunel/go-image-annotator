package label

import (
	"testing"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnLabelUpdateShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	err := repo.Update(lbl.UpdatableModel{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdateLabel(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	name := "a-label"
	label, _ := CreateLabel(repo, name)
	newDesc := "new-description"
	err := repo.Update(lbl.UpdatableModel{Name: label.Name, NewDescription: newDesc})
	assert.Nil(t, err)
	r, err := repo.FindLabel(name)
	assert.Nil(t, err)
	assert.Equal(t, newDesc, r.Description)
}
