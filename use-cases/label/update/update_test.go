package update

import (
	"testing"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.LabelRepo{}, WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := New(&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: non_existing_name}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	repo := &fk.LabelRepo{Existing: []lbl.Label{label}}
	itr := New(repo)
	req := Request{
		Name:           label.Name,
		NewDescription: "updated-description",
	}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, p.Got.Description, req.NewDescription)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	itr := New(&fk.LabelRepo{
		Existing:    []lbl.Label{label},
		ErrOnUpdate: e.ErrInternal,
	})
	itr.Execute(t.Context(), Request{Name: "a-label"}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
