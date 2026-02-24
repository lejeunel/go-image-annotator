package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"modernc.org/sqlite"
	"net/http"
)

var ErrDB = errors.New("database internal error")
var ErrDBForeignKeyConstraint = errors.New("foreign key constraint")
var ErrDBUniqueConstraint = errors.New("unique constraint")
var ErrDBPrimaryKeyConstraint = errors.New("primary key constraint")

var ErrNotFound = errors.New("record not found")
var ErrValidation = errors.New("validation error")
var ErrResourceName = errors.New("only lower case characters and hyphens are allowed: resource name validation error")
var ErrDuplication = errors.New("duplication error")
var ErrDependency = errors.New("dependency error")

var ErrEntitlement = errors.New("entitlement error")
var ErrGroupOwnership = errors.New("group ownership error")
var ErrIdentity = errors.New("identity provider error")
var ErrFileStorage = errors.New("file storage error")

var ErrPagination = errors.New("pagination error")

var ErrForbiddenLabel = errors.New("forbidden annotation error")

type DBErrorConverter interface {
	Convert(error) error
}

type SQLiteErrorConverter struct{}

func (c *SQLiteErrorConverter) Convert(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	var sqliteErr *sqlite.Error
	var errorCode int
	if errors.As(err, &sqliteErr) {
		errorCode = sqliteErr.Code()
		switch errorCode {
		case 2067: //SQLITE_CONSTRAINT_UNIQUE
			return ErrDBUniqueConstraint
		case 1555: //SQLITE_CONSTRAINT_PRIMARYKEY
			return ErrDBPrimaryKeyConstraint
		case 1811: //SQLITE_CONSTRAINT_TRIGGER
			return ErrDBForeignKeyConstraint
		case 787: //SQLITE_CONSTRAINT_TRIGGER
			return ErrDBForeignKeyConstraint
		}
	}
	return fmt.Errorf("%w: %w", ErrDB, err)
}

func ToHTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrDependency):
		return http.StatusFailedDependency
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, ErrDuplication):
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func ToHumaStatusError(err error) huma.StatusError {
	return huma.NewError(ToHTTPStatus(err), err.Error())
}
