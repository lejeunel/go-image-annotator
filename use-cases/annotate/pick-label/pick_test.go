package pick

import (
	"testing"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
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
	labels := []lbl.Label{
		lbl.NewLabel(lbl.NewLabelId(), "fist-label"),
		lbl.NewLabel(lbl.NewLabelId(), "second-label"),
	}
	itr := New(&fk.LabelRepo{Existing: labels}, &fk.ProfileRepo{})
	itr.Execute(t.Context(), nil, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, p.Got[0], labels[0].Name)
}

func TestFetchLabelsInProfile(t *testing.T) {
	p := &FakePresenter{}
	labels := []lbl.Label{
		lbl.NewLabel(lbl.NewLabelId(), "first-label"),
		lbl.NewLabel(lbl.NewLabelId(), "second-label"),
	}
	var labelNames []string
	for _, l := range labels {
		labelNames = append(labelNames, l.Name)
	}
	profile := pr.NewProfile(pr.NewProfileId(), "my-profile", pr.WithLabels(labelNames))
	itr := New(
		&fk.LabelRepo{Existing: labels},
		&fk.ProfileRepo{ExistingProfiles: []pr.Profile{profile}},
	)
	itr.Execute(t.Context(), &profile.Name, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, labelNames, p.Got)
}
