package user

import (
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/add"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/list"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/read"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/update"
)

type Interactors struct {
	Add    add.Interactor
	Delete delete.Interactor
	Update update.Interactor
	List   list.Interactor
	Read   read.Interactor
}
