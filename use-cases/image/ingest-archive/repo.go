package ingest

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

type CollectionRepo interface {
	Find(string) (*clc.Collection, error)
}
