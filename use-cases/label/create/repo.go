package create

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	Create(lbl.Label) error
	Exists(string) (bool, error)
}
