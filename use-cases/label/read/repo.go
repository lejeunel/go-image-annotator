package read

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	FindLabelByName(string) (*l.Label, error)
}
