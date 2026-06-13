package annotation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrieveNilGroupOfAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos()
	labelName := "a-label"
	_, _, _, imageLabel := CreateAnnotedImage(repos, "a-collection", labelName, nil)
	r, err := repos.Annotation.GroupOfAnnotation(imageLabel.Id)
	assert.Nil(t, err)
	assert.Nil(t, r)
}

func TestRetrieveGroupOfAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos()
	labelName := "a-label"
	group := "my-group"
	_, collection, _, imageLabel := CreateAnnotedImage(repos, "a-collection", labelName, &group)
	r, err := repos.Annotation.GroupOfAnnotation(imageLabel.Id)
	assert.Nil(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, (*collection.Group).Name, *r)
}
