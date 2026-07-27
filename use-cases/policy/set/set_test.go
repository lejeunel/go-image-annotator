package set

import (
	"github.com/stretchr/testify/assert"
	"testing"

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
	itr := New(&fk.FileStore{}, a)
	itr.Execute(t.Context(),
		`
version: 1
rules:
    admin:
        - CreateCollection
`,
		p)
	assert.NotNil(t, a.GotRules)
	assert.Contains(t, *a.GotRules, "admin")
}
