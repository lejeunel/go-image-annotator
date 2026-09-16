package sqlite

import (
	"log/slog"

	gr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	lb "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	pr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/profile"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	v "github.com/lejeunel/go-image-annotator/modules/string-validator"
	clc "github.com/lejeunel/go-image-annotator/use-cases/profile"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/create"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/find"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/list"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/update"
)

func NewProfileInteractors(
	pr pr.ProfileRepo,
	gr gr.GroupRepo,
	lr lb.LabelRepo,
	logger slog.Logger,
	pageSize int, auth auth.Interface,
) clc.Interactors {
	return clc.Interactors{
		Find: find.New(pr),
		Create: create.New(pr, lr, gr, create.WithNameValidator(v.NewNameValidator()),
			create.WithAuth(auth)),
		Delete: delete.New(pr, delete.WithAuth(auth)),
		List:   list.New(pr),
		Update: update.New(pr, gr, lr, update.WithAuth(auth)),
	}
}
