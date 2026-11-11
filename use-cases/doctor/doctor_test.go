package doctor

import (
	"testing"

	cfg "github.com/lejeunel/go-image-annotator/config"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	"github.com/stretchr/testify/assert"
)

func TestInvalidAdminPasswordShouldFail(t *testing.T) {
	p := &FakePresenter{}
	pwval := &fk.StringValidator{Invalid: true}
	itr := New(pwval)
	req := Request{cfg.Config{InitialAdminPassword: "invalid-password"}}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.InitialAdminPassword, pwval.Validated)
	assert.Error(t, p.Got.Error)
}

func TestInvalidImageStoreURI(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.StringValidator{})
	req := Request{cfg.Config{ImagesURI: "xyz://invalid-uri"}}
	itr.Execute(t.Context(), req, p)
	assert.Error(t, p.Got.Error)
}
