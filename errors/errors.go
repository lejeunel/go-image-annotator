package errors

import (
	"fmt"
)

type ErrCheckSum struct{}

func (e *ErrCheckSum) Error() string {
	return "Checksum failed."
}

func (e *ErrCheckSum) GetStatus() int {
	return 400
}

type ErrNotFound struct {
	Entity   string
	Criteria string
	Value    string
	Err      error
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("Requested entity of type %s with %s:%s not found, %s", e.Entity, e.Criteria, e.Value, e.Err)
}

func (e *ErrNotFound) GetStatus() int {
	return 404
}

type ErrForbiddenDeletingDependency struct {
	ParentEntity string
	ChildEntity  string
	ParentId     string
}

func (e ErrForbiddenDeletingDependency) Error() string {
	return fmt.Sprintf("Cannot delete %s entity with id %s as it is needed by child resource of type %s",
		e.ParentEntity, e.ParentId, e.ChildEntity)
}

func (e *ErrForbiddenDeletingDependency) GetStatus() int {
	return 403
}

type ErrPermission struct {
	Operation  string
	NeededRole string
	UserRoles  []string
}

func (e ErrPermission) Error() string {
	return fmt.Sprintf("Cannot perform %s. One of %v roles is needed, but you have %v",
		e.Operation, e.NeededRole, e.UserRoles)
}

func (e *ErrPermission) GetStatus() int {
	return 403
}
