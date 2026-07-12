package sqlite

import (
	infra "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	lbl "github.com/lejeunel/go-image-annotator/use-cases/label"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
	"github.com/lejeunel/go-image-annotator/use-cases/label/delete"
	fetchall "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
	"github.com/lejeunel/go-image-annotator/use-cases/label/find"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	"github.com/lejeunel/go-image-annotator/use-cases/label/update"
)

func NewSQLiteLabelInteractors(repo infra.SQLiteLabelRepo, pageSize int, auth auth.Authorizer) lbl.Interactors {
	return lbl.Interactors{
		Find:            *find.New(repo),
		Create:          *create.New(repo, create.WithAuth(auth)),
		Delete:          *delete.New(repo, delete.WithAuth(auth)),
		List:            *list.New(repo),
		Update:          *update.New(repo),
		FetchAll:        *fetchall.New(repo),
		DefaultPageSize: pageSize,
	}
}
