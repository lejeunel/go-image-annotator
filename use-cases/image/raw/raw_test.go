package raw

import (
	"bytes"
	"io"
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrorOnGetRaw(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.FileStore{ErrOnGet: e.ErrNotFound}, &fk.ImageRepo{})
	itr.Execute(im.NewImageId().String(), p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnGetSpecs(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.FileStore{}, &fk.ImageRepo{ErrOnGetSpecs: e.ErrNotFound})
	itr.Execute(im.NewImageId().String(), p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
	assert.False(t, p.GotSuccess)
}

func TestReadRawImage(t *testing.T) {
	p := &FakePresenter{}
	data := []byte("the-data")
	specs := im.ImageSpecs{MIMEType: "the-type"}
	itr := New(&fk.FileStore{Data: data}, &fk.ImageRepo{ReturnSpecs: specs})
	itr.Execute(im.NewImageId().String(), p)
	assert.True(t, p.GotSuccess)
	r, err := io.ReadAll(p.Got)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(data, r))
	assert.Equal(t, specs.MIMEType, p.Got.ImageSpecs.MIMEType)
}
