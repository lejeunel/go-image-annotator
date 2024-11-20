package models

import (
	"github.com/google/uuid"
)

type Set struct {
	Id        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt string    `db:"created_at"`
	UpdatedAt string    `db:"updated_at"`
}
