package profile

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/create"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/find"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/list"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/update"
)

type Interactors struct {
	Find            find.Interactor
	Create          create.Interactor
	Delete          delete.Interactor
	List            list.Interactor
	Update          update.Interactor
	DefaultPageSize int
	Authorizer      auth.Authorizer
}
