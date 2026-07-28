package uow

import (
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	cl "github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
)

type CloneCollectionTransactor struct{ db *sqlx.DB }

func NewCloneCollectionTransactor(db *sqlx.DB) *CloneCollectionTransactor {
	return &CloneCollectionTransactor{db: db}
}

func (u *CloneCollectionTransactor) RunInTx(
	fn func(cl.Repos) error) error {

	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stores := cl.Repos{
		ImageRepo:      im.NewSQLiteImageRepo(tx),
		CollectionRepo: clc.NewSQLiteCollectionRepo(tx),
		AnnotationRepo: an.NewSQLiteAnnotationRepo(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}
	return tx.Commit()
}
