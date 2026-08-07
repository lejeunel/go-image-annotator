package uow

import (
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	cd "github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
)

type DeleteCollectionTransactor struct{ db *sqlx.DB }

func NewDeleteCollectionTransactor(db *sqlx.DB) *DeleteCollectionTransactor {
	return &DeleteCollectionTransactor{db: db}
}

func (u *DeleteCollectionTransactor) RunInTx(
	fn func(cd.Repos) error,
) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stores := cd.Repos{
		ImageRepo:      im.NewSQLiteImageRepo(tx),
		CollectionRepo: clc.NewSQLiteCollectionRepo(tx),
		AnnotationRepo: an.NewSQLiteAnnotationRepo(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}
	return tx.Commit()
}
