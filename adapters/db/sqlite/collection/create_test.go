package collection

import (
	"testing"

	grr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"

	prr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/profile"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewCollectionRepo(db)
	db.Close()
	_, err := CreateCollection(repo, "a-collection")
	assert.ErrorIs(t, err, e.ErrInternal, "expected internal error")
}

func TestCreate(t *testing.T) {
	_, err := CreateCollection(NewCollectionRepo(s.NewInMemory()), "a-collection")
	assert.NoError(t, err, "expected no error on create but got")
}

func TestCreateCollectionInGroup(t *testing.T) {
	db := s.NewInMemory()
	groupRepo := grr.NewGroupRepo(db)
	collectionRepo := NewCollectionRepo(db)
	group := grp.NewGroup(grp.NewGroupId(), "a-group")
	groupRepo.Create(group)
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithGroup(group.Name))
	collectionRepo.Create(c)
	r, err := collectionRepo.GetGroup(c.Name)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "a-group", *r)
}

func TestCreateCollectionWithProfile(t *testing.T) {
	db := s.NewInMemory()
	profileRepo := prr.NewProfileRepo(db)
	collectionRepo := NewCollectionRepo(db)
	profile := pr.NewProfile(pr.NewProfileId(), "my-profile")
	profileRepo.Create(profile)
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithProfile(profile.Name))
	collectionRepo.Create(c)
	r, err := collectionRepo.GetProfile(c.Name)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, profile.Name, *r)
}

func TestCollectionWithoutGroupFailsWithNotFoundErr(t *testing.T) {
	db := s.NewInMemory()
	clcRepo := NewCollectionRepo(db)
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	clcRepo.Create(collection)
	group, _ := clcRepo.GetGroup("a-collection")
	assert.Nil(t, group)
}

func TestNewProfileIsUnused(t *testing.T) {
	db := s.NewInMemory()
	profileRepo := prr.NewProfileRepo(db)
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	profileRepo.Create(profile)

	isUsed, err := profileRepo.IsUsed(profile.Name)
	assert.NoError(t, err)
	assert.NotNil(t, isUsed)
	assert.False(t, *isUsed)
}

func TestAttachProfileToCollection(t *testing.T) {
	db := s.NewInMemory()
	profileRepo := prr.NewProfileRepo(db)
	collectionRepo := NewCollectionRepo(db)

	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	profileRepo.Create(profile)
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithProfile(profile.Name))
	collectionRepo.Create(collection)

	isUsed, err := profileRepo.IsUsed(profile.Name)
	assert.NoError(t, err)
	assert.NotNil(t, isUsed)
	assert.True(t, *isUsed)

	r, err := collectionRepo.Find(collection.Name)
	assert.NoError(t, err)
	assert.NotNil(t, r.Profile)
	assert.Equal(t, profile.Name, *r.Profile)
}
