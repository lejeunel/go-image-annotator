package list

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	List(Request) ([]*l.Label, error)
	Count() (int64, error)
}
