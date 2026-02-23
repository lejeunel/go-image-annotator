package images

import (
	g "datahub/generic"
	sq "github.com/Masterminds/squirrel"
)

func NewBaseImageQuery() sq.SelectBuilder {
	query := g.SqlBuilder.Select(`
				i.id,
				i.camera_id,
				i.captured_at,
				i.created_at,
				i.updated_at,
				i.mimetype,
				i.width,
				i.height,
				i.image_type,
				i.filename,
				i.sha256,
				ic.collection_id AS collection_id,
				c.group_name as group_name`).From("images AS i")
	query = query.LeftJoin("cameras AS ca ON (i.camera_id=ca.id)")
	query = query.Join("image_collection_assoc AS ic ON (ic.image_id=i.id)")
	query = query.Join("collections AS c ON (ic.collection_id=c.id)")

	return query

}

func NewBaseImageCountQuery() sq.SelectBuilder {
	query := g.SqlBuilder.Select("COUNT(*)").From("images AS i")
	query = query.LeftJoin("cameras AS ca ON (i.camera_id=ca.id)")
	query = query.Join("image_collection_assoc AS ic ON (ic.image_id=i.id)")

	return query

}
