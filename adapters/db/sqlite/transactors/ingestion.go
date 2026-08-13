package uow

import (
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	in "github.com/lejeunel/go-image-annotator/modules/image-ingester"
)

type IngestionTransactor struct{ db *sqlx.DB }

func NewIngestionTransactor(
	db *sqlx.DB,
) *IngestionTransactor {
	return &IngestionTransactor{db: db}
}

func (u *IngestionTransactor) RunInTx(
	fn func(in.Repos) error,
) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stores := in.Repos{
		ImageRepo:      im.NewImageRepo(tx, &fk.FilterStrParser{}, &fk.OrderStrParser{}),
		LabelRepo:      lbl.NewLabelRepo(tx),
		CollectionRepo: clc.NewCollectionRepo(tx),
		AnnotationRepo: an.NewAnnotationRepo(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}
	return tx.Commit()
}
