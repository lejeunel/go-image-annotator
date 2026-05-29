package create

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

type Repo interface {
	Create(clc.Collection) error
	Exists(string) (bool, error)
}
