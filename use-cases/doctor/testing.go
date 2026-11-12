package doctor

import (
	"errors"
	"io"
	"strings"

	"github.com/lejeunel/go-image-annotator/config"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
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

type FakePresenter struct {
	LastExam Response
	Errors   []error
}

func (p *FakePresenter) PresentExam(r Response) {
	p.LastExam = r
	if r.Error != nil {
		p.Errors = append(p.Errors, r.Error)
	}
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
