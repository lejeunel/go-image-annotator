package tests

import (
	"bytes"
	"context"
	an "datahub/app/annotator"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	goim "image"
	"math"
	"time"

	a "datahub/app"
	"testing"
)

func InitializeAnnotatorTests(t *testing.T) (context.Context, *a.App, *an.Annotator, *im.Image, *clc.Collection, *lbl.Label) {
	s, _, ctx := a.NewTestApp(t, true)
	label, _ := lbl.New("thelabel", "")
	s.Labels.Create(ctx, label)
	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("thecollection", "", "mygroup")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	ctx = context.WithValue(ctx, "groups", "mygroup")

	annotator := an.NewAnnotator(s.Labels, s.Images,
		s.Collections, s.Locations, s.Authorizer, s.Logger, 640)

	site, _ := loc.NewSite("mysite", "thegroup")
	camera, _ := loc.NewCamera("mycamera", site, "")
	annotator.Locations.SaveSite(ctx, site)
	annotator.Locations.SaveCamera(ctx, camera)

	return ctx, &s, annotator, image, collection, label

}

func TestCreateBoundingBox(t *testing.T) {

	ctx, _, annotator, image, _, label := InitializeAnnotatorTests(t)
	bboxPayload := &an.BoundingBox{Xc: 10, Yc: 10, Height: 4, Width: 3,
		Label: label.Name}
	err := annotator.UpsertBoundingBox(ctx, image.Id, image.CollectionId, bboxPayload)
	AssertNoError(t, err)
	state, err := annotator.MakeState(ctx,
		an.AnnotatorRequest{ImageId: image.Id, CollectionId: image.CollectionId})
	AssertNoError(t, err)
	if len(state.BoundingBoxes) != 1 {
		t.Fatalf("expected to retrieve one bounding box, but got %v", len(state.BoundingBoxes))
	}
}

func TestUpdateBoundingBoxCoords(t *testing.T) {

	ctx, _, annotator, image, _, label := InitializeAnnotatorTests(t)
	bboxPayload := &an.BoundingBox{Xc: 10, Yc: 10, Height: 4, Width: 3,
		Label: label.Name}
	annotator.UpsertBoundingBox(ctx, image.Id, image.CollectionId, bboxPayload)
	state, _ := annotator.MakeState(ctx,
		an.AnnotatorRequest{ImageId: image.Id, CollectionId: image.CollectionId})

	newXValue := bboxPayload.Xc + 1
	bboxPayload.Xc = newXValue
	bboxPayload.Id = state.BoundingBoxes[0].Id
	err := annotator.UpsertBoundingBox(ctx, image.Id, image.CollectionId, bboxPayload)
	AssertNoError(t, err)
	state, _ = annotator.MakeState(ctx,
		an.AnnotatorRequest{ImageId: image.Id, CollectionId: image.CollectionId})

	got := state.BoundingBoxes[0].Xc
	if math.Abs(got-newXValue) > 0.001 {
		t.Fatalf("expected to retrieve updated bbox coordinate (%v), but got %v", newXValue, got)
	}
}

func TestImageResizerResizesWithTargetWidth(t *testing.T) {
	ctx, _, annotator, image, _, _ := InitializeAnnotatorTests(t)

	origWidth := float64(image.Width)

	state, _ := annotator.MakeState(ctx,
		an.AnnotatorRequest{ImageId: image.Id, CollectionId: image.CollectionId})
	transformedImage, _, _ := goim.Decode(bytes.NewReader(state.ImageData))

	expectedWidth := float64(annotator.Rescaler.TargetWidth)
	gotWidth := transformedImage.Bounds().Max.X
	if gotWidth != int(expectedWidth) {
		t.Fatalf("expected image to be resized from %v to %v, but got %v", origWidth, expectedWidth, gotWidth)
	}

}

func TestImageResizerTransformsBBoxCoordinates(t *testing.T) {
	ctx, _, annotator, image, collection, label := InitializeAnnotatorTests(t)
	origWidth := float64(image.Width)
	origHeight := float64(image.Height)
	ratio := origWidth / origHeight
	transformedWidth := float64(annotator.Rescaler.TargetWidth)
	transformedHeight := transformedWidth / ratio

	inputBBox := &an.BoundingBox{Xc: transformedWidth / 2, Yc: transformedHeight / 2,
		Height: transformedHeight, Width: transformedWidth,
		Label: label.Name}
	annotator.UpsertBoundingBox(ctx, image.Id, image.CollectionId, inputBBox)

	retrievedImage, _ := annotator.Images.Find(ctx, image.Id, collection.Id, im.FetchMetaOnly)
	retrievedBBox := retrievedImage.BoundingBoxes[0]

	AssertDeepEqual(t, retrievedBBox.Coords,
		im.BBoxCoords{Xc: origWidth / 2,
			Yc:     origHeight / 2,
			Height: origHeight,
			Width:  origWidth},
		"bounding box coordinates")
}

func TestScrollImages(t *testing.T) {
	s, clock, ctx := a.NewTestApp(t, true)
	label, _ := lbl.New("thelabel", "")
	s.Labels.Create(ctx, label)
	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("thecollection", "", "mygroup")
	s.Collections.Create(ctx, collection)

	ctx = context.WithValue(ctx, "groups", "mygroup")

	scroller := an.NewScroller(s.Images, *im.NewImageFilter(im.WithCollectionId(collection.Id)),
		*im.NewAscendingImageCapturedOrder())
	var images []*im.Image
	for range 3 {
		clock.Advance(1 * time.Hour)
		image, _ = im.New(testJPGImage)
		image.CapturedAt = clock.Now()

		s.Images.Save(ctx, image, collection)
		images = append(images, image)

	}
	firstRightImage, err := scroller.GetNextImage(ctx, images[0])
	AssertNoError(t, err)
	if !images[1].Id.Equal(&firstRightImage.Id) {
		t.Fatalf("expected first right image id to be %v, but got %v",
			images[1].Id, firstRightImage.Id)
	}

	_, err = scroller.GetPrevImage(ctx, images[0])
	AssertError(t, err)

	secondLeftImage, err := scroller.GetPrevImage(ctx, images[1])
	AssertNoError(t, err)
	if !images[0].Id.Equal(&secondLeftImage.Id) {
		t.Fatalf("expected second left image id to be %v, but got %v",
			images[0].Id, secondLeftImage.Id)
	}

	_, err = scroller.GetNextImage(ctx, images[2])
	AssertError(t, err)

}
func TestBoundingBoxAreColorized(t *testing.T) {
	ctx, _, annotator, image, _, label := InitializeAnnotatorTests(t)
	first := &an.BoundingBox{Xc: 10, Yc: 10, Height: 4, Width: 3,
		Label: label.Name}
	annotator.UpsertBoundingBox(ctx, image.Id, image.CollectionId, first)
	second := &an.BoundingBox{Xc: 11, Yc: 11, Height: 5, Width: 4,
		Label: label.Name}
	annotator.UpsertBoundingBox(ctx, image.Id, image.CollectionId, second)
	state, _ := annotator.MakeState(ctx,
		an.AnnotatorRequest{ImageId: image.Id, CollectionId: image.CollectionId})

	firstColor := state.BoundingBoxes[0].Color
	secondColor := state.BoundingBoxes[1].Color
	if firstColor == secondColor {
		t.Fatalf("expected that boxes are assigned different colors, but got %v, and %v",
			firstColor, secondColor)
	}

}
