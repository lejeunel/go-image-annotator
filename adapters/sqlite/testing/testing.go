package testing

import (
	i "github.com/lejeunel/go-image-annotator/adapters/sqlite/infra"
	r "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	f "github.com/lejeunel/go-image-annotator/app/file-store"
)

func NewSQLiteTestingInfra() i.SQLiteInfra {
	db := r.NewSQLiteDB(":memory:")
	return i.NewSQLiteInfra(db, &f.FakeStore{})
}
