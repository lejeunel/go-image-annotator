package metadata

import (
	"testing"

	"github.com/jmoiron/sqlx"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	cr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	ir "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func Init() (*sqlx.DB, MetaRepo, clc.Collection, im.Image) {
	db := s.NewInMemory()
	metaRepo := NewMetaRepo(db)
	clcRepo := cr.NewCollectionRepo(db)
	imRepo := ir.NewImageRepo(db, &fk.FilterStrParser{}, &fk.OrderStrParser{})

	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), collection)
	clcRepo.Create(collection)
	imRepo.AddImage(image.Id, nil, im.Specs{})
	imRepo.AddToCollection(image.Id, collection.Name)

	return db, metaRepo, collection, image
}

func TestInternalErrOnAddShouldFail(t *testing.T) {
	db, repo, collection, image := Init()
	db.Close()
	err := repo.Add(collection.Name, image.Id, "key", "value")
	assert.ErrorIs(t, err, e.ErrInternal, "expected internal error")
}

func TestKeyExists(t *testing.T) {
	_, repo, collection, image := Init()
	err := repo.Add(collection.Name, image.Id, "key", "value")
	assert.NoError(t, err)
	keyExists, err := repo.KeyExists(collection.Name, image.Id, "key")
	assert.NoError(t, err)
	assert.True(t, keyExists)
}

func TestKeyDoesNotExist(t *testing.T) {
	_, repo, collection, image := Init()
	keyExists, err := repo.KeyExists(collection.Name, image.Id, "non-existing-key")
	assert.NoError(t, err)
	assert.False(t, keyExists)
}

func TestGetValue(t *testing.T) {
	_, repo, collection, image := Init()
	repo.Add(collection.Name, image.Id, "key", "value")
	value, err := repo.GetValue(collection.Name, image.Id, "key")
	assert.NoError(t, err)
	assert.Equal(t, "value", *value)
}

func TestUpdateValue(t *testing.T) {
	_, repo, collection, image := Init()
	repo.Add(collection.Name, image.Id, "key", "value")
	err := repo.UpdateValue(collection.Name, image.Id, "key", "new-value")
	assert.NoError(t, err)
	r, _ := repo.GetValue(collection.Name, image.Id, "key")
	assert.Equal(t, "new-value", *r)
}

func TestDeleteOne(t *testing.T) {
	_, repo, collection, image := Init()
	firstKey, firstValue := "first-key", "first-value"
	secondKey, secondValue := "second-key", "second-value"
	repo.Add(collection.Name, image.Id, firstKey, firstValue)
	repo.Add(collection.Name, image.Id, secondKey, secondValue)
	err := repo.Delete(collection.Name, image.Id, firstKey)
	assert.NoError(t, err)
	firstKeyExists, _ := repo.KeyExists(collection.Name, image.Id, firstKey)
	assert.False(t, firstKeyExists)
	secondKeyExists, _ := repo.KeyExists(collection.Name, image.Id, secondKey)
	assert.True(t, secondKeyExists)
}

func TestDeleteAll(t *testing.T) {
	_, repo, collection, image := Init()
	firstKey, firstValue := "first-key", "first-value"
	secondKey, secondValue := "second-key", "second-value"
	repo.Add(collection.Name, image.Id, firstKey, firstValue)
	repo.Add(collection.Name, image.Id, secondKey, secondValue)
	err := repo.DeleteAll(collection.Name, image.Id)
	assert.NoError(t, err)
	firstKeyExists, _ := repo.KeyExists(collection.Name, image.Id, firstKey)
	assert.False(t, firstKeyExists)
	secondKeyExists, _ := repo.KeyExists(collection.Name, image.Id, secondKey)
	assert.False(t, secondKeyExists)
}

func TestList(t *testing.T) {
	_, repo, collection, image := Init()
	firstKey, firstValue := "first-key", "first-value"
	secondKey, secondValue := "second-key", "second-value"
	repo.Add(collection.Name, image.Id, firstKey, firstValue)
	repo.Add(collection.Name, image.Id, secondKey, secondValue)
	meta, err := repo.List(collection.Name, image.Id)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(meta))
	assert.Equal(t, firstValue, meta[0].Value)
}
