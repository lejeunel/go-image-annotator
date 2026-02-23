package generic

import (
	"database/sql/driver"
	"github.com/google/uuid"
)

type UUIDWrapper[T any] struct {
	UUID uuid.UUID
}

// Implement sql.Scanner
func (u *UUIDWrapper[T]) Scan(value interface{}) error {
	var id uuid.UUID
	if err := id.Scan(value); err != nil {
		return err
	}
	u.UUID = id
	return nil
}

// Implement driver.Valuer
func (u UUIDWrapper[T]) Value() (driver.Value, error) {
	return u.UUID.Value()
}

func (u UUIDWrapper[T]) String() string {
	return u.UUID.String()
}
