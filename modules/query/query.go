package query

import (
	"go.tomakado.io/dumbql/query"
)

type SQLizer struct {
	sql  string
	args []any
	err  error
}

func NewSQLizerFromExpr(expr query.Expr) SQLizer {
	sql, args, err := expr.ToSql()
	return SQLizer{sql, args, err}
}

func NewSQLizer(sql string, args []any) SQLizer {
	return SQLizer{sql, args, nil}
}

func (s SQLizer) ToSql() (string, []any, error) {
	return s.sql, s.args, s.err
}
