package find

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	FindLabel(string) (*l.Label, error)
}
