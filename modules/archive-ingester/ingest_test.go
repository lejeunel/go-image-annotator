package ingester

import (
	"archive/zip"
	"bytes"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func MakeZipArchive(files map[string][]byte) (*bytes.Reader, int64) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for name, data := range files {
		w, err := zw.Create(name)
		if err != nil {
			panic(err)
		}

		if _, err := w.Write(data); err != nil {
			panic(err)
		}
	}

	if err := zw.Close(); err != nil {
		panic(err)
	}

	data := buf.Bytes()
	return bytes.NewReader(data), int64(len(data))
}

func Setup() (ArchiveIngester, *bytes.Reader, int64) {
	ing := New(&fk.ImageStore{}, &FakeImageIngester{})
	archive, size := MakeZipArchive(
		map[string][]byte{
			"image1.jpg": st.TestJPGImage,
			"image2.png": st.TestPNGImage,
		},
	)
	return ing, archive, size
}

func TestIngestArchive(t *testing.T) {
	ing, reader, size := Setup()
	r, err := ing.IngestArchive(Request{ReaderAt: reader, Size: size})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(r.ImageIds))
}

func TestDeleteAllOnFailure(t *testing.T) {
	ing, reader, size := Setup()

	s := fk.ImageStore{}
	imageIngester := FakeImageIngester{Err: e.ErrInternal}
	ing.ImageIngester = &imageIngester
	ing.ImageStore = &s
	_, err := ing.IngestArchive(Request{ReaderAt: reader, Size: size})
	assert.Error(t, err)
	assert.True(t, s.DeletedBatch)
}
