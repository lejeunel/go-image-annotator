package label

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	u "github.com/lejeunel/go-image-annotator/use-cases/label/update"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnLabelUpdateShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	err := repo.Update(u.Model{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdateLabel(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	name := "a-label"
	label, _ := CreateLabel(repo, name)
	newDesc := "new-description"
	err := repo.Update(u.Model{Name: label.Name, NewDescription: newDesc})
	assert.Nil(t, err)
	r, err := repo.FindLabel(name)
	assert.Nil(t, err)
	assert.Equal(t, newDesc, r.Description)
}
