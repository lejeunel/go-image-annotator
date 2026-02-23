package labels

import (
	"context"
	g "datahub/generic"
	"io"
)

type LabelDescriber struct {
	LabelService *Service
	Id           LabelId
}

func (c LabelDescriber) Describe(ctx context.Context, w io.Writer) error {

	label, err := c.LabelService.Find(ctx, c.Id)
	if err != nil {
		return err
	}

	description := g.MakeDescriptionTable(map[string]string{
		"id":         label.String(),
		"name":       label.Name,
		"created_at": label.CreatedAt.Format("2006-01-02 / 15:04"),
	},
		[]string{"id", "name", "created_at"})

	description.Render(w)
	return nil
}
