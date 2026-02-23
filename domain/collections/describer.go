package collections

import (
	"context"
	g "datahub/generic"
	"io"
)

type CollectionDescriber struct {
	CollectionService *Service
	Id                CollectionId
}

func (c CollectionDescriber) Describe(ctx context.Context, w io.Writer) error {

	collection, err := c.CollectionService.Find(ctx, c.Id)
	if err != nil {
		return err
	}

	var profileName string
	if collection.Profile == nil {
		profileName = "n/a"
	} else {
		profileName = collection.Profile.Name
	}

	description := g.MakeDescriptionTable(map[string]string{
		"id":                 collection.Id.String(),
		"name":               collection.Name,
		"created_at":         collection.CreatedAt.Format("2006-01-02 / 15:04"),
		"updated_at":         collection.UpdatedAt.Format("2006-01-02 / 15:04"),
		"group":              collection.Group,
		"annotation_profile": profileName,
		"description":        collection.Description},
		[]string{"id", "name", "created_at", "updated_at", "group", "annotation_profile", "description"})

	if err := description.Render(w); err != nil {
		return err
	}
	return nil
}
