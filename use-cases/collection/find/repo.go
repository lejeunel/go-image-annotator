package find

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

type Repo interface {
	Find(string) (*clc.Collection, error)
}
