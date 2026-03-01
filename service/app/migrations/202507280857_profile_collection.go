package migrations

import (
	"context"
	"database/sql"
)

func Up202507300832(ctx context.Context, tx *sql.Tx) error {

	if _, err := tx.Exec(`PRAGMA foreign_keys=off;`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
			CREATE TABLE collections_new (
				id varchar(36) PRIMARY KEY,
                name TEXT UNIQUE,
                created_at timestamp,
                updated_at timestamp,
				description TEXT DEFAULT '',
                group_name TEXT DEFAULT '',
				profile_id VARCHAR(36),
				FOREIGN KEY (profile_id) REFERENCES annotation_profiles(id) ON DELETE SET NULL
			);
		`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
			INSERT INTO collections_new (id, name, created_at, updated_at, description, group_name, profile_id)
			SELECT id, name, created_at, updated_at, description, group_name, NULL FROM collections;
		`); err != nil {
		return err
	}

	if _, err := tx.Exec(`DROP TABLE collections;`); err != nil {
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE collections_new RENAME TO collections;`); err != nil {
		return err
	}

	if _, err := tx.Exec(`PRAGMA foreign_keys=on;`); err != nil {
		return err
	}

	return nil
}

func Down202507300832(ctx context.Context, tx *sql.Tx) error {

	if _, err := tx.Exec(`PRAGMA foreign_keys=off;`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
			CREATE TABLE collections_old (
				id varchar(36) PRIMARY KEY,
                name TEXT UNIQUE,
				description TEXT DEFAULT ''
                created_at timestamp,
                updated_at timestamp,
				description TEXT DEFAULT '',
                group_name TEXT DEFAULT ''
			);
		`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
			INSERT INTO collections_old (id, name, description, created_at, updated_at, group_name)
			SELECT id, name, description, created_at, updated_at, group_name FROM collections;
		`); err != nil {
		return err
	}

	if _, err := tx.Exec(`DROP TABLE collections;`); err != nil {
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE collections_old RENAME TO collections;`); err != nil {
		return err
	}

	if _, err := tx.Exec(`PRAGMA foreign_keys=on;`); err != nil {
		return err
	}

	return nil
}
