package label

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/stretchr/testify/assert"
)

func TestUpdateLabel(t *testing.T) {
	repo := NewSQLiteLabelRepo(s.NewInMemory())
	name := "a-label"
	label, _ := CreateLabel(repo, name)
	newDesc := "new-description"
	err := repo.Update(lbl.UpdatableModel{Name: label.Name, NewDescription: newDesc})
	assert.Nil(t, err)
	r, err := repo.FindLabel(name)
	assert.Nil(t, err)
	assert.Equal(t, newDesc, r.Description)
}
