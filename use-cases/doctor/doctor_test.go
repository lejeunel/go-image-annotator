package doctor

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	"github.com/stretchr/testify/assert"
)

func TestInvalidAdminPasswordShouldFail(t *testing.T) {
	p := &FakePresenter{}
	pwval := &fk.StringValidator{Invalid: true}
	itr := New(pwval)
	req := Request{InitialAdminPassword: "invalid-password"}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.InitialAdminPassword, pwval.Validated)
	assert.Error(t, p.Got.Exams[InitialAdminPassword].Error)
}
