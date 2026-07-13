package role

import (
	"github.com/lejeunel/go-image-annotator/use-cases/role/create"
	"github.com/lejeunel/go-image-annotator/use-cases/role/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/role/find"
	"github.com/lejeunel/go-image-annotator/use-cases/role/list"
	"github.com/lejeunel/go-image-annotator/use-cases/role/update"
)

type Interactors struct {
	Find            find.Interactor
	Create          create.Interactor
	Delete          delete.Interactor
	List            list.Interactor
	Update          update.Interactor
	DefaultPageSize int
}
