package find

import (
	"github.com/stretchr/testify/assert"
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestAuth(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.FileStore{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotAuthErr)
}

func TestStoreErr(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.FileStore{ErrOnGet: e.ErrInternal})
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestReadPolicy(t *testing.T) {
	p := &FakePresenter{}
	data := []byte("test-data")
	itr := New(&fk.FileStore{Data: data})
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, data, p.Got)
}
