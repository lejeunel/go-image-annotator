package db

import (
	"database/sql"
)

type Querier interface {
	Get(dest any, query string, args ...any) error
	Select(dest any, query string, args ...any) error
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}
