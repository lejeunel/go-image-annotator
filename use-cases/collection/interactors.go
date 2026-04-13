package collection

import (
	"github.com/lejeunel/go-image-annotator-v2/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/collection/read"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/collection/update"
)

type Interactors struct {
	Find            read.Interactor
	Create          create.Interactor
	Delete          delete.Interactor
	List            list.Interactor
	Update          update.Interactor
	DefaultPageSize int
}
