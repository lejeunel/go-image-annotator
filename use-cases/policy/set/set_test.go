package set

import (
	"testing"

	"github.com/stretchr/testify/assert"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestAuth(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.FileStore{}, &fk.Auth{Err: e.ErrAuthorization})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
}

func TestInvalidYamlForm(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.FileStore{}, &fk.Auth{})
	itr.Execute(t.Context(), "invalid-payload", p)
	assert.True(t, p.GotValidationErr)
}

func TestInvalidMethod(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.FileStore{}, &fk.Auth{})
	itr.Execute(t.Context(),
		`
version: 1
rules:
    admin:
        - NonExistinMethod
`,
		p)
	assert.True(t, p.GotValidationErr)
}

func TestSetNewRules(t *testing.T) {
	a := &fk.Auth{}
	p := &FakePresenter{}
	fs := &fk.FileStore{}
	itr := New(fs, a)
	data := `
version: 1
rules:
    admin:
        - CreateCollection
`
	itr.Execute(t.Context(), data, p)
	assert.NotNil(t, a.GotRules)
	assert.Contains(t, *a.GotRules, "admin")
	assert.Equal(t, string(fs.GotData), data)
}
