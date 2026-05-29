package validation

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
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
			if !tt.isValid && !errors.Is(err, e.ErrValidation) {
				t.Fatalf("expected validation error given name %v, got none", tt.name)
			}
			if tt.isValid && err != nil {
				t.Fatalf("expected no validation error given name %v, got one", tt.name)
			}
		})
	}

}
