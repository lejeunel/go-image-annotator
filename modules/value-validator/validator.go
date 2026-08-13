package value_validator

import (
	"fmt"
	"time"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Validator interface {
	Validate(any) error
}

type BaseTypeValidator struct{}

func (v BaseTypeValidator) Validate(value any) error {
	if value == "" {
		return fmt.Errorf("string value cannot be empty: %w", e.ErrValidation)
	}

	switch value.(type) {
	case int, float32, float64, string, time.Time, bool:
		return nil
	case nil:
		return fmt.Errorf("value cannot be nil: %w", e.ErrValidation)
	default:
		return fmt.Errorf("unsupported value type: %T: %w", v, e.ErrValidation)
	}
}
