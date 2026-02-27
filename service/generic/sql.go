package generic

import (
	sq "github.com/Masterminds/squirrel"
)

var SqlBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
