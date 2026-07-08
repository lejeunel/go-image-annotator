package raw

import (
	"bytes"
	"io"
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrorOnFind(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeFileGetter{Err: e.ErrNotFound})
	itr.Execute(Request{ImageId: im.NewImageId().String()}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
	assert.False(t, p.GotSuccess)
}

func TestReadRawImage(t *testing.T) {
	p := &FakePresenter{}
	data := []byte("the-data")
	itr := New(&FakeFileGetter{data: data})
	itr.Execute(Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotSuccess)
	r, err := io.ReadAll(p.Got)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(data, r))
}
