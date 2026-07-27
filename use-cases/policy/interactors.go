package user

import (
	"github.com/lejeunel/go-image-annotator/use-cases/policy/read"
	"github.com/lejeunel/go-image-annotator/use-cases/policy/set"
)

type Interactors struct {
	Read read.Interactor
	Set  set.Interactor
}
