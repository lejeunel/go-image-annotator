package collection

import (
	"time"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

func CreateCollection(repo *SQLiteCollectionRepo, name string) (*clc.Collection, error) {
	c := clc.NewCollection(clc.NewCollectionId(), name,
		clc.WithDescription("a-description"), clc.WithCreatedAt(time.Now()),
		clc.WithGroup("a-group"))
	if err := repo.Create(c); err != nil {
		return nil, err
	}
	return &c, nil

}
