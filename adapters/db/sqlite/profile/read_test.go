package profile

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewProfileRepo(s.NewInMemory())
	CreateProfile(repo, "a-profile", nil)
	_, err := repo.Find("non-existing-profile")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	CreateProfile(repo, "a-profile", nil)
	db.Close()
	_, err := repo.Find("a-profile")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRetrieve(t *testing.T) {
	repo := NewProfileRepo(s.NewInMemory())
	c := pr.NewProfile(pr.NewProfileId(), "a-profile",
		pr.WithDescription("a-description"))
	repo.Create(c)
	r, err := repo.Find("a-profile")
	assert.NoError(t, err, "expected no error on find")
	assert.Equal(t, c.Name, r.Name)
	assert.Equal(t, c.Description, r.Description)
	assert.Equal(t, c.Id, r.Id)
}
