package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	e "datahub/errors"
	g "datahub/generic"
	"testing"
)

func InitializeAnnotationTests(t *testing.T) (*a.App, *clc.Collection, *im.Image, *lbl.Label, context.Context) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "admin")
	ctx = context.WithValue(ctx, "groups", "mygroup")
	image, _ := im.New(testPNGImage)
	collection, _ := clc.New("myimageset", "", "mygroup")
	label, _ := lbl.New("my-label", "mydescription")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)
	s.Labels.Create(ctx, label)

	return &s, collection, image, label, ctx
}

func TestListImageOfLabel(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationTests(t)

	firstLabel, _ := lbl.New("firstlabel", "")
	secondLabel, _ := lbl.New("secondlabel", "")
	s.Labels.Create(ctx, firstLabel)
	s.Labels.Create(ctx, secondLabel)

	firstImage, _ := im.New(testJPGImage)
	secondImage, _ := im.New(testJPGImage)
	collection, _ := clc.New("thecollection", "", "")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, firstImage, collection)
	s.Images.Save(ctx, secondImage, collection)

	s.Images.Annotations.ApplyLabel(ctx, firstLabel, firstImage)
	s.Images.Annotations.ApplyLabel(ctx, secondLabel, secondImage)

	retrieved, _, err := s.Images.List(ctx, im.FilterArgs{LabelId: &secondLabel.Id},
		im.OrderingArgs{}, g.PaginationParams{Page: 1, PageSize: 2}, im.FetchMetaOnly)
	AssertNoError(t, err)

	if len(retrieved) != 1 {
		t.Fatalf("expected to retrieve 1 image, but got %v", len(retrieved))
	}
	retrievedLabel := retrieved[0].Annotations[0].Label
	if retrievedLabel.Id != secondLabel.Id {
		t.Fatalf("expected to retrieve image with label %v, but got %v",
			secondLabel.String(), retrievedLabel.String())
	}

}

func TestApplyingDuplicateLabelShouldFail(t *testing.T) {
	s, _, image, label, ctx := InitializeAnnotationTests(t)
	err := s.Images.Annotations.ApplyLabel(ctx, label, image)
	AssertNoError(t, err)

	err = s.Images.Annotations.ApplyLabel(ctx, label, image)
	AssertErrorIs(t, err, e.ErrDuplication)

}

func TestDeleteLabeledImageAndItsAssociatedLabel(t *testing.T) {
	s, _, image, label, ctx := InitializeAnnotationTests(t)

	s.Images.Annotations.ApplyLabel(ctx, label, image)
	s.Images.RemoveFromCollection(ctx, image)
	s.Labels.Delete(ctx, label)
	res, _ := s.Labels.Find(ctx, label.Id)
	if res != nil {
		t.Fatal("expected to retrieve 0 labels, but got one")
	}

}

func TestRemovingLabelFromImage(t *testing.T) {
	s, collection, image, label, ctx := InitializeAnnotationTests(t)
	err := s.Images.Annotations.ApplyLabel(ctx, label, image)
	err = s.Images.Annotations.Delete(ctx, image.Annotations[0].Id.String())
	AssertNoError(t, err)

	retrievedImages, _, err := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		*im.NewImageDefaultOrderingArgs(),
		g.OneItemPaginationParams,
		im.FetchMetaOnly)

	if len(retrievedImages[0].Annotations) != 0 {
		t.Fatalf("expected to retrieve image with no annotations, but got %v",
			len(retrievedImages[0].Annotations))
	}

}
