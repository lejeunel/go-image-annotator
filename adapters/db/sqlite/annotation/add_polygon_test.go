package annotation

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

var TestingPolygonPoints = a.Points{Coordinates: [][2]float32{{0, 0}, {1, 1}}}

func TestInternalErrOnAddPolygonShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	polygon := a.NewPolygon(a.NewAnnotationId(), TestingPolygonPoints, label)
	repos.Annotation.Db.Close()
	err := repos.Annotation.AddPolygon(image.Id, collection.Id, polygon, nil, nil)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnFindPolygonShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	polygon := a.NewPolygon(a.NewAnnotationId(), TestingPolygonPoints, label)
	repos.Annotation.AddPolygon(image.Id, collection.Id, polygon, nil, nil)
	repos.Annotation.Db.Close()
	_, err := repos.Annotation.FindPolygons(image.Id, collection.Id)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddPolygon(t *testing.T) {
	repos := NewAnnotationTestRepos()
	labelName := "a-label"
	image, collection, label := CreateAnnotableImage(repos, "a-collection", labelName, nil)
	polygon := a.NewPolygon(a.NewAnnotationId(), TestingPolygonPoints, label)
	err := repos.Annotation.AddPolygon(image.Id, collection.Id, polygon, nil, nil)
	assert.NoError(t, err)
	r, err := repos.Annotation.FindPolygons(image.Id, collection.Id)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(r))
	assert.Equal(t, labelName, r[0].Label.Name)
	assert.Equal(t, polygon.Points, r[0].Points)
}
