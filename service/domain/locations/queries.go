package locations

import (
	g "datahub/generic"
	sq "github.com/Masterminds/squirrel"
)

func NewBaseSiteQuery() sq.SelectBuilder {
	return g.SqlBuilder.Select(`
			DISTINCT
				s.id,
				s.name,
				s.group_name,
				s.created_at,
				s.updated_at
			`).From("sites AS s")
}
