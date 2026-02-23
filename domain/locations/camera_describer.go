package locations

import (
	"context"
	g "datahub/generic"
	"io"
)

type CameraDescriber struct {
	LocationService *Service
	Id              CameraId
}

func (c CameraDescriber) Describe(ctx context.Context, w io.Writer) error {

	camera, err := c.LocationService.FindCamera(ctx, c.Id)
	if err != nil {
		return err
	}

	site, err := c.LocationService.FindSite(ctx, camera.Site.Id)
	if err != nil {
		return err
	}

	description := g.MakeDescriptionTable(map[string]string{
		"id":          camera.Id.String(),
		"name":        camera.Name,
		"site_name":   site.Name,
		"site_id":     site.Id.String(),
		"transmitter": camera.Transmitter,
		"created_at":  camera.CreatedAt.Format("2006-01-02 / 15:04"),
	},
		[]string{"id", "name", "site_name", "site_id", "transmitter", "created_at"})

	description.Render(w)
	return nil
}
