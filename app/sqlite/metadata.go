package sqlite

import (
	cr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	ir "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	mr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/metadata"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	kv "github.com/lejeunel/go-image-annotator/modules/string-validator"
	vv "github.com/lejeunel/go-image-annotator/modules/value-validator"
	mu "github.com/lejeunel/go-image-annotator/use-cases/metadata"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/add"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/list"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/update"
)

func NewSQLiteMetadataInteractors(
	mr mr.SQLiteMetaRepo,
	cr cr.SQLiteCollectionRepo,
	ir ir.SQLiteImageRepo,
	auth auth.Interface) mu.Interactors {
	return mu.Interactors{
		Add:    add.New(cr, ir, mr, kv.NewNameValidator(), vv.BaseTypeValidator{}, add.WithAuth(auth)),
		Delete: delete.New(cr, mr, delete.WithAuth(auth)),
		Update: update.New(cr, ir, mr, update.WithAuth(auth)),
		List:   list.New(mr),
	}
}
