package interactors

import (
	infra "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/label"
	lbl "github.com/lejeunel/go-image-annotator/use-cases/label"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
	"github.com/lejeunel/go-image-annotator/use-cases/label/delete"
	fetchall "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	"github.com/lejeunel/go-image-annotator/use-cases/label/read"
)

func NewSQLiteLabelInteractors(repo *infra.SQLiteLabelRepo, pageSize int) *lbl.Interactors {
	return &lbl.Interactors{
		Find:            *read.New(repo),
		Create:          *create.New(repo),
		Delete:          *delete.New(repo),
		List:            *list.New(repo),
		FetchAll:        *fetchall.New(repo),
		DefaultPageSize: pageSize,
	}
}
