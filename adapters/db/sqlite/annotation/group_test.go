package annotation

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNilGroupOfAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos(s.NewInMemory())
	labelName := "a-label"
	_, _, _, imageLabel := CreateAnnotedImage(repos, "a-collection", labelName, nil)
	r, err := repos.Annotation.GroupOfAnnotation(imageLabel.Id)
	assert.Nil(t, err)
	assert.Nil(t, r)
}

func TestRetrieveGroupOfAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos(s.NewInMemory())
	labelName := "a-label"
	group := "my-group"
	_, collection, _, imageLabel := CreateAnnotedImage(repos, "a-collection", labelName, &group)
	r, err := repos.Annotation.GroupOfAnnotation(imageLabel.Id)
	assert.Nil(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, (*collection.Group).Name, *r)
}
