package sql

type SQLizer interface {
	ToSql() (string, []any, error)
}
