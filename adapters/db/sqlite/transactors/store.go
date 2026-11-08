package uow

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	m "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/metadata"
	s "github.com/lejeunel/go-image-annotator/modules/image-store"
	qu "github.com/lejeunel/go-image-annotator/modules/query"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type StoreTransactor struct {
	db *sqlx.DB
	qu.FilterSQLizer
	im.OrderStrParser
}

func NewStoreTransactor(
	db *sqlx.DB,
	fp qu.FilterSQLizer,
	op im.OrderStrParser,
) *StoreTransactor {
	return &StoreTransactor{db, fp, op}
}

func (u *StoreTransactor) RunInTx(
	fn func(s.Repos) error,
) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning transaction: %v: %w", err, e.ErrInternal)
	}
	defer tx.Rollback()

	stores := s.Repos{
		ImageRepo:      im.NewImageRepo(tx, u.FilterSQLizer, u.OrderStrParser),
		CollectionRepo: clc.NewCollectionRepo(tx),
		AnnotationRepo: an.NewAnnotationRepo(tx),
		MetaRepo:       m.NewMetaRepo(tx),
	}

	if err := fn(stores); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %v: %w", err, e.ErrInternal)
	}
	return nil
}
