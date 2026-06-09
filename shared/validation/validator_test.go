package validation

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var nameTests = []struct {
	name    string
	isValid bool
}{
	{"name", true},
	{"a-name", true},
	{"NAME", false},
	{"%^&*()", false},
	{"a-name-0001", true},
}

func TestNameValidator(t *testing.T) {
	for _, tt := range nameTests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewNameValidator()
			err := validator.Validate(tt.name)
			if !tt.isValid {
				assert.ErrorIs(t, err, e.ErrValidation)
			}
			if tt.isValid {
				assert.NoError(t, err)
			}
		})
	}

}
