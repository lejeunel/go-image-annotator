package profile

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnProfileCountShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	_, err := repo.Count()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCountProfiles(t *testing.T) {
	repo := NewProfileRepo(s.NewInMemory())
	CreateProfile(repo, "a-profile")
	count, _ := repo.Count()
	assert.Equal(t, 1, int(*count))
}

func TestInternalErrOnProfileListShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewProfileRepo(db)
	db.Close()
	_, err := repo.List(pa.PaginationParams{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestListProfiles(t *testing.T) {
	repo := NewProfileRepo(s.NewInMemory())
	CreateProfile(repo, "a-profile")
	CreateProfile(repo, "another-profile")
	cs, err := repo.List(pa.PaginationParams{Page: 1, PageSize: 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cs))
	assert.False(t, cs[0].Name == cs[1].Name)
}
