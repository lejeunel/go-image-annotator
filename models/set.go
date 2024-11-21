package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"regexp"
)

type Set struct {
	Id        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt string    `db:"created_at"`
	UpdatedAt string    `db:"updated_at"`
}

func (s Set) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Name, validation.Match(regexp.MustCompile("^[a-z\\-_]*$"))),
	)
}
