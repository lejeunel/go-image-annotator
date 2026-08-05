package value_validator

import (
	"testing"
	"time"
)

func TestValidateValue(t *testing.T) {
	validator := BaseTypeValidator{}
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"int", 42, false},
		{"float32", float32(3.14), false},
		{"float64", 3.14, false},
		{"string", "hello", false},
		{"bool", true, false},
		{"time", time.Now(), false},

		{"nil", nil, true},
		{"slice", []int{1, 2, 3}, true},
		{"map", map[string]string{"a": "b"}, true},
		{"struct", struct{}{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
