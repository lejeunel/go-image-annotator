package label

import (
	"github.com/lejeunel/go-image-annotator-v2/use-cases/label/create"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/label/delete"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/label/list"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/label/read"
)

type Interactors struct {
	Find            read.Interactor
	Create          create.Interactor
	Delete          delete.Interactor
	List            list.Interactor
	DefaultPageSize int
}
