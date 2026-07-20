package uow

import (
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	in "github.com/lejeunel/go-image-annotator/modules/ingester"
)

type IngestionUoW struct{ db *sqlx.DB }

func NewIngestionUoW(db *sqlx.DB) *IngestionUoW { return &IngestionUoW{db: db} }

func (u *IngestionUoW) RunInTx(
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
