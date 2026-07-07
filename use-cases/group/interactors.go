package group

import (
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"github.com/lejeunel/go-image-annotator/use-cases/group/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/group/find"
	"github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

type Interactors struct {
	Find            find.Interactor
	Create          create.Interactor
	Delete          delete.Interactor
	List            list.Interactor
	Update          update.Interactor
	DefaultPageSize int
}
