package doctor

import (
	"errors"

	"github.com/lejeunel/go-image-annotator/config"
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
