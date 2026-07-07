package label

import (
	"github.com/lejeunel/go-image-annotator/modules/auth"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
	"github.com/lejeunel/go-image-annotator/use-cases/label/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
	"github.com/lejeunel/go-image-annotator/use-cases/label/find"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
)

type Interactors struct {
	Find            find.Interactor
	Create          create.Interactor
	Delete          delete.Interactor
	List            list.Interactor
	FetchAll        fetchall.Interactor
	DefaultPageSize int
	Authorizer      auth.Auth
}
