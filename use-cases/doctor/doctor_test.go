package doctor

import (
	"testing"

	cfg "github.com/lejeunel/go-image-annotator/config"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func NewTestingInteractor() Interactor {
	return New(FakeConfigProvider{},
		&FakePasswordValidatorProvider{&fk.StringValidator{}})
}

func TestFetchingConfigFailure(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ConfigProvider = FakeConfigProvider{Err: e.ErrConfigParsing}
	itr.Execute(t.Context(), p)
	assert.Error(t, p.LastExam.Error)
}

func TestInvalidAdminPasswordShouldFail(t *testing.T) {
	p := &FakePresenter{}
	pwval := &fk.StringValidator{Invalid: true}
	itr := NewTestingInteractor()
	cfg := cfg.Config{InitialAdminPassword: "invalid-password"}
	itr.ConfigProvider = FakeConfigProvider{Config: cfg}
	itr.PasswordValidatorProvider = &FakePasswordValidatorProvider{pwval}
	itr.Execute(t.Context(), p)
	assert.Equal(t, cfg.InitialAdminPassword, pwval.Validated)
	assert.True(t, p.GotErr(e.ErrValidation))
}

func TestInvalidImageStoreURI(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	cfg := cfg.Config{ImagesURI: "xyz://invalid-uri"}
	itr.ConfigProvider = FakeConfigProvider{Config: cfg}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrValidation))
}
