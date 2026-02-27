package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	e "datahub/errors"
	g "datahub/generic"
	clk "github.com/jonboulle/clockwork"
	"testing"
	"time"
)

func InitializeBoundingBoxTests(t *testing.T) (context.Context, *a.App, *im.Image, *clc.Collection, *lbl.Label, *clk.FakeClock) {
	s, clock, ctx := a.NewTestApp(t, true)
	label, _ := lbl.New("thelabel", "")
	s.Labels.Create(ctx, label)
	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("thecollection", clc.WithGroup("mygroup"))
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	ctx = context.WithValue(ctx, "groups", "mygroup")

	return ctx, &s, image, collection, label, clock
}

func TestBoundingBoxesCoordinates(t *testing.T) {
	tests := map[string]struct {
		xc float64
		yc float64
		w  float64
		h  float64
		ok bool
	}{
		"negative x coord should fail":       {xc: -2, yc: 4, h: 3, w: 5, ok: false},
		"negative y coord should fail":       {xc: 2, yc: -4, h: 3, w: 5, ok: false},
		"all positive coords should succeed": {xc: 2, yc: 4, h: 3, w: 5, ok: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := im.NewBoundingBox(tc.xc, tc.yc, tc.h, tc.w)
			if !tc.ok {
				AssertErrorIs(t, err, e.ErrValidation)
			}

		})
	}

}

func TestApplyingBoundingBoxesToImage(t *testing.T) {
	ctx, s, image, collection, label, _ := InitializeBoundingBoxTests(t)
	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)

	bbox.Annotate(label)
	err := s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertNoError(t, err)

	retrievedImage, err := s.Images.Find(ctx, image.Id, collection.Id, im.FetchMetaOnly)
	AssertNoError(t, err)

	retrievedBoxes := retrievedImage.BoundingBoxes
	if len(retrievedBoxes) != 1 {
		t.Fatalf("expected to retrieve one bounding box, but got %v", len(retrievedBoxes))
	}

	if retrievedBoxes[0].Coords != bbox.Coords {
		t.Fatalf("expected to retrieve bounding box coordinates %v, but got %v",
			bbox.Coords, retrievedBoxes[0].Coords)
	}
}

func TestDeleteBoundingBoxById(t *testing.T) {
	ctx, s, image, collection, label, _ := InitializeBoundingBoxTests(t)
	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)

	err := s.Images.Annotations.Delete(ctx, bbox.Annotation.Id.String())
	AssertNoError(t, err)

	retrievedImage, _ := s.Images.Find(ctx, image.Id, collection.Id, im.FetchMetaOnly)

	if len(retrievedImage.BoundingBoxes) > 0 {
		t.Fatalf("expected to retrieve no annotation, but got some.")
	}

}

func TestUpsertBoundingBoxToImage(t *testing.T) {
	ctx, s, image, collection, label, _ := InitializeBoundingBoxTests(t)

	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	if len(image.BoundingBoxes) < 1 {
		t.Fatal("expected to retrieve 1 bounding-box after upsert, but got none")
	}

	bbox.Coords.Xc = 0
	bbox.Coords.Yc = 9
	err := s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertNoError(t, err)

	retrievedImage, _ := s.Images.Find(ctx, image.Id, collection.Id, im.FetchMetaOnly)
	if len(retrievedImage.BoundingBoxes) != 1 {
		t.Fatalf("expected to retrieve 1 bounding box, got %v", len(retrievedImage.BoundingBoxes))
	}
	AssertDeepEqual(t, bbox.Coords, retrievedImage.BoundingBoxes[0].Coords, "bounding_box")
}

func TestListImagesByLabel(t *testing.T) {
	ctx, s, image, collection, label, _ := InitializeBoundingBoxTests(t)
	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)

	bbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)

	notAnnotatedImage, _ := im.New(testPNGImage)
	s.Images.Save(ctx, notAnnotatedImage, collection)

	retrievedImages, _, _ := s.Images.List(ctx, *im.NewImageFilter(im.WithLabelId(label.Id)),
		im.OrderingArgs{}, g.PaginationParams{Page: 1, PageSize: 2}, im.FetchMetaOnly)
	if len(retrievedImages) != 1 {
		t.Fatalf("expected to retrieve one image, but got %v", len(retrievedImages))
	}

}

func TestListImagesByLabelWithTwoImages(t *testing.T) {
	ctx, s, image, collection, label, _ := InitializeBoundingBoxTests(t)
	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)

	otherImage, _ := im.New(testPNGImage)
	s.Images.Save(ctx, otherImage, collection)
	otherBbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	otherBbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, otherBbox, otherImage)

	retrievedImages, _, _ := s.Images.List(ctx, *im.NewImageFilter(im.WithLabelId(label.Id)),
		im.OrderingArgs{}, g.PaginationParams{Page: 1, PageSize: 3}, im.FetchMetaOnly)
	if len(retrievedImages) != 2 {
		t.Fatalf("expected to retrieve two images, but got %v", len(retrievedImages))
	}

}

func TestUpdateLabel(t *testing.T) {
	ctx, s, image, collection, label, _ := InitializeBoundingBoxTests(t)

	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)

	newLabel, _ := lbl.New("newlabel", "")
	err := s.Labels.Create(ctx, newLabel)
	AssertNoError(t, err)
	err = s.Images.Annotations.UpdateLabel(ctx, bbox.Annotation.Id.String(), "newlabel")
	AssertNoError(t, err)

	retrievedImage, _ := s.Images.Find(ctx, image.Id, collection.Id, im.FetchMetaOnly)
	name := retrievedImage.BoundingBoxes[0].Annotation.Label.Name
	if name != "newlabel" {
		t.Fatalf("expected to retrieve %v name, got %v", "newlabel", name)
	}

}

func TestApplyingBBoxWithoutLabelShouldFail(t *testing.T) {
	ctx, s, image, _, _, _ := InitializeBoundingBoxTests(t)
	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)

	err := s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertErrorIs(t, err, e.ErrValidation)

}

func TestBoundingBoxesAreReturnedInAntiChronologicalOrder(t *testing.T) {
	ctx, s, image, _, label, clock := InitializeBoundingBoxTests(t)
	firstBbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	firstBbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, firstBbox, image)

	clock.Advance(1 * time.Hour)
	secondBbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	secondBbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, secondBbox, image)

	image, _ = s.Images.Find(ctx, image.Id, image.Collection.Id, im.FetchMetaOnly)
	if len(image.BoundingBoxes) != 2 {
		t.Fatalf("expected to retrieve 2 bounding boxes but got %v", len(image.BoundingBoxes))
	} else {
		got := image.BoundingBoxes[0].Annotation.CreatedAt
		want := secondBbox.Annotation.CreatedAt
		if got.Equal(want) == false {
			t.Fatalf("expected to retrieve last annotation at first index (timestamp: %v), but got %v",
				want, got)
		}
	}

}
