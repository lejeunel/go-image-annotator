package pick

import (
	"slices"
	"testing"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnCount: e.ErrInternal}, &fk.ProfileRepo{})
	itr.Execute(t.Context(), nil, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrWhenCountExceedsLimit(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{Count_: 2}, &fk.ProfileRepo{}, WithLimit(1))
	itr.Execute(t.Context(), nil, p)
	assert.ErrorIs(t, p.GotErr, e.ErrLabelLimitExceeded)
}

func TestHandleErrOnFetch(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnFetch: e.ErrInternal}, &fk.ProfileRepo{})
	itr.Execute(t.Context(), nil, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestFetchLabels(t *testing.T) {
	p := &FakePresenter{}
	labels := []string{"first-label", "second-labels"}
	itr := New(&fk.LabelRepo{ExistingNames: labels}, &fk.ProfileRepo{})
	itr.Execute(t.Context(), nil, p)
	assert.True(t, p.GotSuccess)
	assert.True(t, slices.Equal(p.Got, labels))
}

func TestFetchLabelsInProfile(t *testing.T) {
	p := &FakePresenter{}
	labels := []string{"first-label", "second-labels"}
	profile := pr.NewProfile(pr.NewProfileId(), "my-profile", pr.WithLabels(labels))
	itr := New(&fk.LabelRepo{ExistingNames: labels}, &fk.ProfileRepo{Return: profile})
	itr.Execute(t.Context(), &profile.Name, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, labels, p.Got)
}
