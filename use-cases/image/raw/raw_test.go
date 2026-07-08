package raw

import (
	"bytes"
	"io"
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrorOnGetRaw(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeFileGetter{Err: e.ErrNotFound}, &FakeRepo{})
	itr.Execute(im.NewImageId().String(), p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnGetSpecs(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeFileGetter{}, &FakeRepo{Err: e.ErrNotFound})
	itr.Execute(im.NewImageId().String(), p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
	assert.False(t, p.GotSuccess)
}

func TestReadRawImage(t *testing.T) {
	p := &FakePresenter{}
	data := []byte("the-data")
	specs := &im.ImageSpecs{MIMEType: "the-type"}
	itr := New(&FakeFileGetter{data: data}, &FakeRepo{ReturnSpecs: specs})
	itr.Execute(im.NewImageId().String(), p)
	assert.True(t, p.GotSuccess)
	r, err := io.ReadAll(p.Got)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(data, r))
	assert.Equal(t, specs.MIMEType, p.Got.ImageSpecs.MIMEType)
}
