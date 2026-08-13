package uow

import (
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	m "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/metadata"
	s "github.com/lejeunel/go-image-annotator/modules/image-store"
)

type StoreTransactor struct {
	db *sqlx.DB
	im.FilterStrParser
	im.OrderStrParser
}

func NewStoreTransactor(
	db *sqlx.DB,
	fp im.FilterStrParser,
	op im.OrderStrParser,
) *StoreTransactor {
	return &StoreTransactor{db, fp, op}
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
		ImageRepo:      im.NewImageRepo(tx, u.FilterStrParser, u.OrderStrParser),
		CollectionRepo: clc.NewCollectionRepo(tx),
		AnnotationRepo: an.NewAnnotationRepo(tx),
		MetaRepo:       m.NewMetaRepo(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}
	return tx.Commit()
}
