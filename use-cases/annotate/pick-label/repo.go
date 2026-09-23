package pick

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type LabelRepo interface {
	FetchAll() ([]string, error)
	Count() (int64, error)
}

type ProfileRepo interface {
	Find(string) (*pr.Profile, error)
}
