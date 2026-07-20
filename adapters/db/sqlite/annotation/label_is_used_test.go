package annotation

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnLabelIsUsedShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repos := NewAnnotationTestRepos(db)
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	imLabel := a.NewImageLabel(label)
	repos.Annotation.AddImageLabel(image.Id, collection.Id, imLabel, nil, nil)
	db.Close()
	_, err := repos.Label.IsUsed(label.Name)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestLabelIsUsedbyAnnotation(t *testing.T) {
	db := s.NewInMemory()
	repos := NewAnnotationTestRepos(db)
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	imLabel := a.NewImageLabel(label)
	repos.Annotation.AddImageLabel(image.Id, collection.Id, imLabel, nil, nil)
	isUsed, err := repos.Label.IsUsed(label.Name)
	assert.NoError(t, err)
	assert.True(t, *isUsed)
}
