package ingest

import (
	"bytes"
	"context"
	"io"

	ast "github.com/lejeunel/go-image-annotator/app/file-store"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

func NewTestingInteractor(opts ...Option) *Interactor {
	i := &Interactor{
		ImageRepo:          &FakeImageRepo{},
		CollectionRepo:     &FakeCollectionRepo{},
		LabelRepo:          &FakeLabelRepo{},
		AnnotationRepo:     &FakeAnnotationRepo{},
		ArtefactRepo:       &ast.FakeStore{},
		Hasher:             &FakeHasher{},
		Logger:             logging.NewNoOpLogger(),
		ImageSpecsDetector: &FakeSpecsDetector{},
		auth:               auth.PassThroughAuth{},
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

type FakeHasher struct {
	sum []byte
}

func (f *FakeHasher) Write(p []byte) (int, error) {
	return len(p), nil
}

func (f *FakeHasher) Sum(b []byte) []byte {
	return append(b, f.sum...)
}

func (f *FakeHasher) Reset() {}

func (f *FakeHasher) Size() int {
	return len(f.sum)
}

func (f *FakeHasher) BlockSize() int {
	return 1
}

type FakePresenter struct {
	Got        *Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(r Response) {
	p.Got = &r
	p.GotSuccess = true
}

type FakeCollectionRepo struct {
	Err                 error
	ErrOnFindCollection bool
	MissingCollection   bool
}

type FakeLabelRepo struct {
	Err              error
	ErrOnLabelExists bool
	MissingLabel     bool
}

type FakeAnnotationRepo struct {
	Err                   error
	ErrOnAddBoundingBox   bool
	ErrOnAddLabel         bool
	NumLabelsAdded        int
	NumBoundingboxesAdded int
}

type FakeImageRepo struct {
	Err                  error
	GotImage             bool
	GotHash              []byte
	GotSpecs             im.ImageSpecs
	ErrOnAddToCollection bool
	ErrOnAddImage        bool
	ErrOnFindHash        bool
	ErrOnDeleteImage     bool
	HashAlreadyExists    bool
	NumDeletedImages     int
}

func (r *FakeCollectionRepo) FindCollectionByName(name string) (*clc.Collection, error) {
	if r.ErrOnFindCollection {
		return nil, r.Err
	}
	if r.MissingCollection {
		return nil, e.ErrNotFound
	}
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithGroup(g.NewGroup(g.NewGroupId(), "my-group")))
	return &c, nil
}

func (r *FakeLabelRepo) FindLabel(name string) (*lbl.Label, error) {
	if r.ErrOnLabelExists {
		return nil, r.Err
	}
	if r.MissingLabel {
		return nil, e.ErrNotFound
	}
	l := lbl.NewLabel(lbl.NewLabelId(), name)
	return &l, nil
}

func (r *FakeImageRepo) FindImageIdByHash(hash []byte) (*im.ImageId, error) {
	if r.ErrOnFindHash {
		return nil, r.Err
	}
	if r.HashAlreadyExists {
		existingId := im.NewImageId()
		return &existingId, nil
	}
	return nil, e.ErrNotFound
}

func (r *FakeAnnotationRepo) AddImageLabel(im.ImageId, clc.CollectionId, an.ImageLabel) error {
	if r.ErrOnAddLabel {
		return r.Err
	}
	r.NumLabelsAdded += 1
	return nil
}

func (r *FakeAnnotationRepo) AddBoundingBox(im.ImageId, clc.CollectionId, an.BoundingBox) error {
	if r.ErrOnAddBoundingBox {
		return r.Err
	}
	r.NumBoundingboxesAdded += 1
	return nil
}

func (r *FakeImageRepo) Delete(im.ImageId) error {
	if r.ErrOnDeleteImage {
		return r.Err
	}
	r.NumDeletedImages += 1
	return nil
}

func (r *FakeImageRepo) AddToCollection(im.ImageId, clc.CollectionId) error {
	if r.ErrOnAddToCollection {
		return r.Err
	}
	return nil
}

func (r *FakeImageRepo) AddImage(imageId im.ImageId, hash []byte, specs im.ImageSpecs) error {
	if r.ErrOnAddImage {
		return r.Err
	}
	r.GotHash = hash
	r.GotSpecs = specs
	return nil
}

type FakeImageReader struct {
	Buffer bytes.Buffer
	Err    error
}

func (d *FakeImageReader) Read(b []byte) (int, error) {
	if d.Err != nil {
		return 0, d.Err
	}
	return d.Buffer.Read(b)

}

type FakeSpecsDetector struct {
	Err    error
	Return im.ImageSpecs
}

func (d *FakeSpecsDetector) Detect(r io.Reader) (*im.ImageSpecs, io.Reader, error) {
	if d.Err != nil {
		return nil, nil, d.Err
	}
	return &d.Return, r, nil

}

type FailingAuth struct {
}

func (f FailingAuth) IngestImage(ctx context.Context, group string) error {
	return e.ErrAuth
}
