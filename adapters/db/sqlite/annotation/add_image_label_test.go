package annotation

import (
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnAddLabelShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repos := NewAnnotationTestRepos(db)
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	db.Close()
	err := repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label), nil, nil)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnFindImageLabels(t *testing.T) {
	db := s.NewInMemory()
	repos := NewAnnotationTestRepos(db)
	image, collection, _ := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	db.Close()
	_, err := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddAndRetrieveImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos(s.NewInMemory())
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
