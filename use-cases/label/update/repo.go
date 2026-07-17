package update

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	Update(lbl.UpdatableModel) error
	Exists(string) (bool, error)
}
