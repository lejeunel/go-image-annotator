package sqlite

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	pl "github.com/lejeunel/go-image-annotator/use-cases/policy"
	"github.com/lejeunel/go-image-annotator/use-cases/policy/read"
)

func NewSQLitePolicyInteractors(fs read.Store, auth auth.Interface) pl.Interactors {
	return pl.Interactors{
		Read: read.New(fs, read.WithAuth(auth)),
	}
}
