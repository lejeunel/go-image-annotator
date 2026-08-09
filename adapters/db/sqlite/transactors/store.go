package uow

import (
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	m "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/metadata"
	s "github.com/lejeunel/go-image-annotator/modules/image-store"
)

type StoreTransactor struct{ db *sqlx.DB }

func NewStoreTransactor(
	db *sqlx.DB,
) *StoreTransactor {
	return &StoreTransactor{db: db}
}

func (u *StoreTransactor) RunInTx(
	fn func(s.Repos) error,
) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stores := s.Repos{
		ImageRepo:      im.NewSQLiteImageRepo(tx),
		CollectionRepo: clc.NewSQLiteCollectionRepo(tx),
		AnnotationRepo: an.NewSQLiteAnnotationRepo(tx),
		MetaRepo:       m.NewSQLiteMetaRepo(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}
	return tx.Commit()
}
