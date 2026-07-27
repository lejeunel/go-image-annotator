package sqlite

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pl "github.com/lejeunel/go-image-annotator/use-cases/policy"
	"github.com/lejeunel/go-image-annotator/use-cases/policy/read"
	"github.com/lejeunel/go-image-annotator/use-cases/policy/set"
)

func NewSQLitePolicyInteractors(fs fs.Interface, auth auth.Interface) pl.Interactors {
	return pl.Interactors{
		Read: read.New(fs, read.WithAuth(auth)),
		Set:  set.New(fs, auth),
	}
}
