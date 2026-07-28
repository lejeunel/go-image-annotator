package uow

import (
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	in "github.com/lejeunel/go-image-annotator/modules/ingester"
	cd "github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
)

type IngestionTransactor struct{ db *sqlx.DB }

func NewIngestionTransactor(db *sqlx.DB) *IngestionTransactor { return &IngestionTransactor{db: db} }

func (u *IngestionTransactor) RunInTx(
	fn func(in.Repos) error) error {

	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stores := in.Repos{
		ImageRepo:      im.NewSQLiteImageRepo(tx),
		LabelRepo:      lbl.NewSQLiteLabelRepo(tx),
		CollectionRepo: clc.NewSQLiteCollectionRepo(tx),
		AnnotationRepo: an.NewSQLiteAnnotationRepo(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}
	return tx.Commit()
}

type DeleteCollectionTransactor struct{ db *sqlx.DB }

func NewDeleteCollectionTransactor(db *sqlx.DB) *DeleteCollectionTransactor {
	return &DeleteCollectionTransactor{db: db}
}

func (u *DeleteCollectionTransactor) RunInTx(
	fn func(cd.Repos) error) error {

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
