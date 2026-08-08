package metadata

import (
	"database/sql"
	"encoding/json"
	"fmt"

	adb "github.com/lejeunel/go-image-annotator/adapters/db"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type SQLiteMetaRepo struct {
	Db adb.Querier
}

type Row struct {
	ImageId      im.ImageId       `db:"image_id"`
	CollectionId clc.CollectionId `db:"collection_id"`
	Meta         json.RawMessage  `db:"meta"`
}

func (r SQLiteMetaRepo) Add(
	collection clc.CollectionName,
	imageId im.ImageId,
	key string,
	value any,
) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = r.Db.Exec(`
        INSERT INTO metadata (image_id, collection_id, meta)
        VALUES (
			?,
			(SELECT id FROM collections WHERE name=?),
			json_object(?, json(?)))
        ON CONFLICT(image_id, collection_id) DO UPDATE
        SET meta = json_set(
            metadata.meta,
            ?,
            json(?)
        )
    `,
		imageId,
		collection,
		key,
		string(valueJSON),
		"$."+key,
		string(valueJSON),
	)
	if err != nil {
		return fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	return err
}

func (r SQLiteMetaRepo) KeyExists(
	collection clc.CollectionName,
	imageId im.ImageId,
	key string,
) (bool, error) {
	var exists bool

	err := r.Db.Get(
		&exists,
		`
        SELECT json_type(meta, ?) IS NOT NULL
        FROM metadata
        WHERE image_id = ?
          AND collection_id = (
              SELECT id
              FROM collections
              WHERE name = ?
          )
        `,
		"$."+key,
		imageId,
		collection,
	)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	return exists, nil
}

func (r SQLiteMetaRepo) GetValue(
	collection clc.CollectionName,
	imageID im.ImageId,
	key string,
) (*any, error) {
	var value any

	err := r.Db.Get(
		&value,
		`
        SELECT json_extract(meta, ?)
        FROM metadata
        WHERE image_id = ?
          AND collection_id = (
              SELECT id
              FROM collections
              WHERE name = ?
          )
        `,
		"$."+key,
		imageID,
		collection,
	)
	if err == sql.ErrNoRows {
		return nil, e.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	return &value, nil
}

func (r SQLiteMetaRepo) UpdateValue(
	collection clc.CollectionName,
	imageID im.ImageId,
	key string,
	value any,
) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	res, err := r.Db.Exec(`
        UPDATE metadata
        SET meta = json_set(meta, ?, json(?))
        WHERE image_id = ?
          AND collection_id = (
              SELECT id
              FROM collections
              WHERE name = ?
          )
    `,
		"$."+key,
		string(valueJSON),
		imageID,
		collection,
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%v: %w", err, e.ErrInternal)
	}
	if n == 0 {
		return fmt.Errorf("no rows affected: %w", e.ErrValidation)
	}

	return nil
}

func (r SQLiteMetaRepo) Delete(
	collection clc.CollectionName,
	imageID im.ImageId,
	key string,
) error {
	res, err := r.Db.Exec(`
        UPDATE metadata
        SET meta = json_remove(meta, ?)
        WHERE image_id = ?
          AND collection_id = (
              SELECT id
              FROM collections
              WHERE name = ?
          )
    `,
		"$."+key,
		imageID,
		collection,
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%v: %w", err, e.ErrInternal)
	}
	if n == 0 {
		return fmt.Errorf("no rows affected: %w", e.ErrValidation)
	}

	return nil
}

func (r SQLiteMetaRepo) List(
	collection clc.CollectionName,
	imageId im.ImageId,
) ([]m.MetaData, error) {
	var metadata []m.MetaData

	if err := r.Db.Select(&metadata,
		`
        SELECT
            je.key,
            je.value
        FROM metadata m,
             json_each(m.meta) AS je
        WHERE m.image_id = ?
          AND m.collection_id = (
              SELECT id
              FROM collections
              WHERE name = ?
          )
    `, imageId, collection); err != nil {
		return nil, fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	return metadata, nil
}

func NewSQLiteMetaRepo(db adb.Querier) SQLiteMetaRepo {
	return SQLiteMetaRepo{Db: db}
}
