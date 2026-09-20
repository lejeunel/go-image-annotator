package create

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := "my-group"
	itr := New(&fk.CollectionRepo{},
		&fk.GroupRepo{Return: g.NewGroup(g.NewGroupId(), "a-group")},
		&fk.ProfileRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Group: &group}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateCollectionWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-collection"
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ErrOnCreate: e.ErrDuplicate},
		&fk.GroupRepo{},
		&fk.ProfileRepo{})
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateCollectionWithInvalidNameShouldFail(t *testing.T) {
	name := "my-collection%/"
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ErrOnCreate: e.ErrValidation},
		&fk.GroupRepo{},
		&fk.ProfileRepo{},
		WithNameValidator(&fk.StringValidator{Invalid: true}))
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotValidationErr)
}

func TestCreateCollectionInNonExistingGroupShouldFail(t *testing.T) {
	name := "my-collection"
	group := "non-existing-group"
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{},
		&fk.GroupRepo{ErrOnFind: e.ErrNotFound},
		&fk.ProfileRepo{},
	)
	itr.Execute(t.Context(), Request{Name: name, Group: &group}, p)
	assert.True(t, p.GotNotFoundErr)
}

func TestCreateCollectionWithNonExistingProfileShouldFail(t *testing.T) {
	name := "my-collection"
	profile := "non-existing-profile"
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{},
		&fk.GroupRepo{},
		&fk.ProfileRepo{ErrOnFind: e.ErrNotFound},
	)
	itr.Execute(t.Context(), Request{Name: name, Profile: &profile}, p)
	assert.True(t, p.GotNotFoundErr)
}

func TestCreateCollection(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.CollectionRepo{}
	now := time.Now()
	group := g.NewGroup(g.NewGroupId(), "my-group")
	profile := pr.NewProfile(pr.NewProfileId(), "my-profile")
	itr := New(repo,
		&fk.GroupRepo{Return: group},
		&fk.ProfileRepo{Return: profile},
		WithClock(clockwork.NewFakeClockAt(now)))
	req := Request{
		Name:        "a-name",
		Description: "a-description",
		Group:       &group.Name,
		Profile:     &profile.Name,
	}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.Name, repo.Got.Name)
	assert.Equal(t, *req.Group, *repo.Got.Group)
	assert.Equal(t, *req.Profile, *repo.Got.Profile)
	assert.Equal(t, req.Description, repo.Got.Description)
	assert.Equal(t, now, repo.Got.CreatedAt)
	assert.False(t, repo.Got.Id.IsNil())
}
