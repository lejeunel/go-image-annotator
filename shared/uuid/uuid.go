package uuid

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type UUIDWrapper[T any] struct {
	UUID uuid.UUID
}

func FromUUID[T any](id uuid.UUID) UUIDWrapper[T] {
	return UUIDWrapper[T]{UUID: id}
}

// Implement sql.Scanner
func (u *UUIDWrapper[T]) Scan(value any) error {
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

func (u UUIDWrapper[T]) IsNil() bool {
	return u.UUID == uuid.Nil
}

func ShortenUUID(id string) string {
	return id[:8]
}
