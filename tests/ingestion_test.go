package tests

import (
	"context"
	app "datahub/app"
	in "datahub/app/ingester"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	e "datahub/errors"
	g "datahub/generic"
	base64 "encoding/base64"
	"testing"
)

func InitializeIngestionTest(t *testing.T) (*app.App, *clc.Collection, *loc.Site, *loc.Camera, *lbl.Label, context.Context) {
	a, _, ctx := app.NewTestApp(t, true)

	ctx = context.WithValue(ctx, "groups", "my-group")
	ctx = context.WithValue(ctx, "entitlements", "im-contrib|annotation-contrib")

	collection, _ := clc.New("my-collection", "collection-description", "my-group")
	site, _ := loc.NewSite("my-site", "my-group")
	camera, _ := loc.NewCamera("my-camera", site, "")
	label, _ := lbl.New("my-label", "")
	a.Collections.Create(ctx, collection)
	a.Locations.SaveSite(ctx, site)
	err := a.Locations.SaveCamera(ctx, camera)
	AssertNoError(t, err)
	a.Labels.Create(ctx, label)

	return &a, collection, site, camera, label, ctx

}

func TestIngestionWithUnspecifiedLocationShouldFail(t *testing.T) {

	base64Data := base64.StdEncoding.EncodeToString(testJPGImage)
	payload := in.ImageIngestionPayload{
		Data:       base64Data,
		Site:       nil,
		Camera:     nil,
		CapturedAt: "2024-05-10T11:04:05.000Z",
		Group:      "my-group",
	}

	a, collection, _, _, _, ctx := InitializeIngestionTest(t)
	ctx = context.WithValue(ctx, "groups", "my-group")

	_, err := a.Ingestion.Ingest(ctx, collection.Name, payload)
	AssertNoError(t, err)
	images, _, _ := a.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		*im.NewImageDefaultOrderingArgs(),
		g.OneItemPaginationParams,
		im.FetchWithRawData)
	if len(images) != 1 {
		t.Fatalf("expected to ingest one image but got %v", len(images))
	}

	if images[0].Camera != nil {
		t.Fatal("expected to ingest image without location, but got one")
	}

}

func TestIngestingImageWithBoundingBox(t *testing.T) {
	a, collection, _, _, _, ctx := InitializeIngestionTest(t)
	ctx = context.WithValue(ctx, "groups", "my-group")

	base64Data := base64.StdEncoding.EncodeToString(testJPGImage)
	payload := in.ImageIngestionPayload{
		Data:       base64Data,
		Site:       nil,
		Camera:     nil,
		CapturedAt: "2024-05-10T11:04:05.000Z",
		Group:      "my-group",
	}

	bbox, _ := im.NewBoundingBox(10, 20, 3, 6)
	bboxPayLoad := in.BoundingBoxIngestion{Xc: bbox.Coords.Xc,
		Yc: bbox.Coords.Yc, Width: bbox.Coords.Width, Height: bbox.Coords.Height,
		Angle: bbox.Coords.Angle, Label: "my-label"}
	payload.BoundingBoxes = []in.BoundingBoxIngestion{bboxPayLoad}

	_, err := a.Ingestion.Ingest(ctx, collection.Name, payload)
	AssertNoError(t, err)
	images, _, _ := a.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		*im.NewImageDefaultOrderingArgs(),
		g.OneItemPaginationParams,
		im.FetchWithRawData)
	if len(images) != 1 {
		t.Fatalf("expected to ingest one image but got %v", len(images))
	}
	if len(images[0].BoundingBoxes) != 1 {
		t.Fatalf("expected to ingest an image with one bounding box but got %v",
			len(images[0].BoundingBoxes))
	}

}

func TestIngestingTwoImagesWithSameByteContentShouldFail(t *testing.T) {
	a, collection, _, _, _, ctx := InitializeIngestionTest(t)
	image, _ := im.New(testJPGImage)
	b64Data := base64.StdEncoding.EncodeToString(image.Data)

	payload := in.ImageIngestionPayload{
		Data:       b64Data,
		CapturedAt: "2024-05-10T11:04:05.000Z",
		Group:      "my-group",
	}
	_, err := a.Ingestion.Ingest(ctx, collection.Name, payload)
	AssertNoError(t, err)

	_, err = a.Ingestion.Ingest(ctx, collection.Name, payload)
	AssertErrorIs(t, err, e.ErrDuplication)

}

func TestIngestingTwoImagesWithSameByteContentShouldSucceedIfOneIsDeleted(t *testing.T) {
	a, collection, site, camera, _, ctx := InitializeIngestionTest(t)
	image, _ := im.New(testJPGImage)
	b64Data := base64.StdEncoding.EncodeToString(image.Data)

	payload := in.ImageIngestionPayload{
		Data:       b64Data,
		Site:       &site.Name,
		Camera:     &camera.Name,
		CapturedAt: "2024-05-10T11:04:05.000Z",
		Group:      "my-group",
	}
	_, err := a.Ingestion.Ingest(ctx, collection.Name, payload)
	AssertNoError(t, err)

	images, _, _ := a.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionName(collection.Name)),
		*im.NewImageDefaultOrderingArgs(), g.OneItemPaginationParams,
		im.FetchMetaOnly)

	a.Images.RemoveFromCollection(ctx, &images[0])

	_, err = a.Ingestion.Ingest(ctx, collection.Name, payload)
	AssertNoError(t, err)

}

func TestIngestingPayloadWithDryRunFlagShouldNotAddImage(t *testing.T) {
	a, collection, site, camera, _, ctx := InitializeIngestionTest(t)
	image, _ := im.New(testJPGImage)
	b64Data := base64.StdEncoding.EncodeToString(image.Data)

	payload := in.ImageIngestionPayload{
		Data:       b64Data,
		Site:       &site.Name,
		Camera:     &camera.Name,
		CapturedAt: "2024-05-10T11:04:05.000Z",
		Group:      "my-group",
		DryRun:     true,
	}
	a.Ingestion.Ingest(ctx, collection.Name, payload)
	images, _, _ := a.Images.List(ctx, im.FilterArgs{},
		im.OrderingArgs{}, g.PaginationParams{Page: 1, PageSize: 1}, im.FetchMetaOnly)
	if len(images) > 0 {
		t.Fatal("expected to find no images when dry-run is set, but got some")
	}

}
