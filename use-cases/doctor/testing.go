package doctor

import (
	"errors"
	"io"
	"strings"

	"github.com/lejeunel/go-image-annotator/config"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type FakeConfigProvider struct {
	Err    error
	Config config.Config
}

func (p FakeConfigProvider) Get() (*config.Config, error) {
	if p.Err != nil {
		return nil, p.Err
	}
	return &p.Config, nil
}

type FakePasswordValidatorProvider struct {
	validator pv.PasswordValidator
}

func (p *FakePasswordValidatorProvider) Get(minEntropy int) pv.PasswordValidator {
	return p.validator
}

type FakeFileStoreProvider struct {
	fs.FileStore
	Err error
}

func (p *FakeFileStoreProvider) Get(config.Config) (fs.FileStore, error) {
	if p.Err != nil {
		return nil, p.Err
	}
	return p.FileStore, nil
}

// FakeDBStoreProvider returns FileStore, or Err when set, and records the
// path it was given.
type FakeDBStoreProvider struct {
	FileStore fs.FileStore
	Err       error
	GotPath   string
}

func (p *FakeDBStoreProvider) Get(path string) (fs.FileStore, error) {
	p.GotPath = path
	if p.Err != nil {
		return nil, p.Err
	}
	return p.FileStore, nil
}

type FakePresenter struct {
	LastExam      Response
	Errors        []error
	CheckedImages int
	// CheckedImagesAtDone is CheckedImages when DoneCheckImageIntegrity was
	// last called, or -1 if it was never called.
	CheckedImagesAtDone int
	NumDone             int
}

func NewFakePresenter() *FakePresenter {
	return &FakePresenter{CheckedImagesAtDone: -1}
}

func (p *FakePresenter) PresentExam(r Response) {
	p.LastExam = r
	if r.Error != nil {
		p.Errors = append(p.Errors, r.Error)
	}
}

func (p *FakePresenter) SuccessCheckImageIntegrity(id im.ImageId) {
	p.CheckedImages += 1
}

func (p *FakePresenter) ErrorCheckImageIntegrity(id im.ImageId) {
	p.CheckedImages += 1
	p.Errors = append(p.Errors, e.ErrNotFound)
}

func (p *FakePresenter) DoneCheckImageIntegrity() {
	p.NumDone += 1
	p.CheckedImagesAtDone = p.CheckedImages
	p.LastExam = Response{}
}

func (p *FakePresenter) GotErr(target error) bool {
	for _, got := range p.Errors {
		if errors.Is(got, target) {
			return true
		}
	}
	return false
}

// CorruptedFileStore returns content other than what was stored.
type CorruptedFileStore struct {
	*fk.FileStore
}

func (s CorruptedFileStore) Get(string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("corrupted")), nil
}

type FakeImageStoreProvider struct {
	ImageStore
}

func (p FakeImageStoreProvider) Get(config.Config) (ImageStore, error) {
	return p.ImageStore, nil
}
