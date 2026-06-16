package annotation

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInternalErrOnAddLabelShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	repos.Annotation.Db.Close()
	err := repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label), nil, nil)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnFindImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, _ := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	repos.Annotation.Db.Close()
	_, err := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddAndRetrieveImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)

	user := u.NewUser("user@example.com")
	repos.User.Create(user)
	now := time.Now()
	repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label), &user.Id, &now)
	labels, err := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(labels))
	assert.NotNil(t, labels[0].Author)
	assert.Equal(t, user.Id, *labels[0].Author)
	assert.NotNil(t, labels[0].Time)
}
