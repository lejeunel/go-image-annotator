package ingest

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

type CollectionRepo interface {
	FindCollectionByName(string) (*clc.Collection, error)
}
