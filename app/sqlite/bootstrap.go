package sqlite

import (
	rr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	ur "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
	bst "github.com/lejeunel/go-image-annotator/use-cases/bootstrap"
)

func NewBootstrapInteractor(userRepo ur.UserRepo, roleRepo rr.RoleRepo,
	fileStore fs.FileStore,
	t tk.TokenHasher, pv pw.PasswordValidator,
) bst.Interactor {
	return bst.New(userRepo, roleRepo, fileStore, t, pv)
}
