package value_validator

import (
	"fmt"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"time"
)

type Validator interface {
	Validate(any) error
}

type BaseTypeValidator struct{}

func (v BaseTypeValidator) Validate(value any) error {
	switch value.(type) {
	case int, float32, float64, string, time.Time, bool:
		return nil
	case nil:
		return fmt.Errorf("value cannot be nil: %w", e.ErrValidation)
	default:
		return fmt.Errorf("unsupported value type: %T: %w", v, e.ErrValidation)
	}
}
