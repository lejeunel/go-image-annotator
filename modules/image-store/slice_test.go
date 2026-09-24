package image_store

import (
	"testing"

	_ "github.com/lejeunel/go-image-annotator/fakes"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrOnSlice(t *testing.T) {
	store, _, image, _ := Setup()
	repos := SetupRepos(image)
	repos.ImageRepo = &fk.ImageRepo{ErrOnSlice: e.ErrInternal}
	store.Repos = repos
	query := "collection:my-collection"
	_, _, err := store.Slice(query, "", pa.PaginationParams{Page: 1, PageSize: 1})
	assert.Error(t, err)
}
