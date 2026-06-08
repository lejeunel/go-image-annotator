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
		Find:            *read.NewInteractor(repo),
		Create:          *create.NewInteractor(repo),
		Delete:          *delete.NewInteractor(repo),
		List:            *list.NewInteractor(repo),
		FetchAll:        *fetchall.NewInteractor(repo),
		DefaultPageSize: pageSize,
	}
}
