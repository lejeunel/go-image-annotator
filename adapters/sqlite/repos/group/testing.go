package group

import (
	g "github.com/lejeunel/go-image-annotator/entities/group"
)

func CreateGroup(repo *SQLiteGroupRepo, name string) (*g.Group, error) {
	c := g.NewGroup(g.NewGroupId(), name,
		g.WithDescription("a-description"))

	if err := repo.Create(c); err != nil {
		return nil, err
	}
	return &c, nil

}
