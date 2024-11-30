package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"regexp"
	"time"
)

type Label struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (l Label) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Name, validation.Required),
		validation.Field(&l.Name, validation.Match(regexp.MustCompile("^[a-z\\-]*$"))),
	)
}
