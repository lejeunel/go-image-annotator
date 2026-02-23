package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	e "datahub/errors"
	g "datahub/generic"
	"fmt"
	"testing"
	"time"

	clk "github.com/jonboulle/clockwork"
)

func InitializeCollectionTests(t *testing.T) (*a.App, *clc.Collection, *im.Image, *clk.FakeClock, context.Context) {
	s, clock, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "admin")
	ctx = context.WithValue(ctx, "groups", "mygroup")
	collection, _ := clc.New("mycollection", "", "mygroup")
	s.Collections.Create(ctx, collection)
	image, _ := im.New(testPNGImage)
	s.Images.Save(ctx, image, collection)

	return &s, collection, image, clock, ctx

}

func TestCreateCollection(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	clc, _ := clc.New("new-collection", "thedescription", "mygroup")
	err := s.Collections.Create(ctx, clc)
	AssertNoError(t, err)

	res, err := s.Collections.Find(ctx, clc.Id)
	AssertNoError(t, err)

	AssertDeepEqual(t, res, clc, "collection")
}

func TestRetrieveCollectionByName(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	clc, _ := clc.New("new-collection", "thedescription", "mygroup")
	err := s.Collections.Create(ctx, clc)
	AssertNoError(t, err)

	retrieved, err := s.Collections.FindByName(ctx, "new-collection")
	AssertNoError(t, err)
	AssertDeepEqual(t, clc, retrieved, "collection")

}

func TestUpdateCollection(t *testing.T) {
	s, clock, ctx := a.NewTestApp(t, true)
	ctx = context.WithValue(ctx, "groups", "mygroup|myothergroup")

	collection, _ := clc.New("mycollection", "thedescription", "mygroup")
	s.Collections.Create(ctx, collection)

	clock.Advance(1 * time.Hour)

	_, err := s.Collections.Update(ctx, collection.Name,
		clc.CollectionUpdatables{Name: "newname", Description: "newdescription",
			Group: "myothergroup"})
	AssertNoError(t, err)

	res, _ := s.Collections.Find(ctx, collection.Id)

	if res.Name != "newname" {
		t.Fatalf("expected to retrieve updated name %v, but got %v",
			"newname", res.Name)
	}
	if res.Description != "newdescription" {
		t.Fatalf("expected to retrieve updated description %v, but got %v",
			"newdescription", res.Description)
	}

	if !res.UpdatedAt.After(collection.UpdatedAt) {
		t.Fatalf("expected to retrived updated timestamp. Got original %v and retrieved %v",
			collection.UpdatedAt, res.UpdatedAt)
	}
	if res.Group != "myothergroup" {
		t.Fatalf("expected to retrieve updated group %v, but got %v",
			"myothergroup", res.Group)
	}

}

func TestCollectionDuplicateNameShouldFail(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	collection, _ := clc.New("mycollection", "", "")
	newCollection, _ := clc.New("mycollection", "", "")
	s.Collections.Create(ctx, collection)
	err := s.Collections.Create(ctx, newCollection)

	AssertErrorIs(t, err, e.ErrDuplication)
}

func TestValidationCollectionName(t *testing.T) {
	tests := map[string]struct {
		name    string
		isValid bool
	}{
		"empty should fail":               {name: "", isValid: false},
		"spaces should fail":              {name: "my set", isValid: false},
		"uppercase should fail":           {name: "MySet", isValid: false},
		"specials should fail":            {name: "my&^/set", isValid: false},
		"spaces and specials should fail": {name: "my &*set", isValid: false},
		"dash should succeed":             {name: "my-set", isValid: true},
		"underscore should succeed":       {name: "my_set", isValid: true},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {
			_, err := clc.New(tc.name, "", "mygroup")
			if tc.isValid {
				AssertNoError(t, err)
			} else {
				AssertErrorIs(t, err, e.ErrResourceName)
			}
		})
	}

}
func TestGetCollectionByName(t *testing.T) {
	s, collection, _, _, ctx := InitializeCollectionTests(t)
	res, err := s.Collections.FindByName(ctx, collection.Name)
	AssertNoError(t, err)
	if res.Name != collection.Name {
		t.Fatalf("expected to retrieve collection with name %v, but got %v",
			collection.Name, res.Name)
	}

}

func TestRetrieveImagesOfCollectionByName(t *testing.T) {

	s, collection, _, _, ctx := InitializeCollectionTests(t)

	page, _, err := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionName(collection.Name)),
		im.OrderingArgs{},
		g.PaginationParams{Page: 1, PageSize: 4},
		im.FetchMetaOnly)
	AssertNoError(t, err)

	if len(page) != 1 {
		t.Fatalf("expected to retrieve 1 image in collection %v, but got %v",
			collection.Name, len(page))
	}

}

func TestRetrieveImagesOfCollectionById(t *testing.T) {

	s, _, ctx := a.NewTestApp(t, true)
	collectionName := "myset"
	collection, _ := clc.New(collectionName, "", "mygroup")
	err := s.Collections.Create(ctx, collection)
	AssertNoError(t, err)
	image, _ := im.New(testPNGImage)
	s.Images.Save(ctx, image, collection)

	page, _, err := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		im.OrderingArgs{},
		g.PaginationParams{Page: 1, PageSize: 4},
		im.FetchMetaOnly)
	AssertNoError(t, err)

	if len(page) != 1 {
		t.Fatalf("expected to retrieve 1 image in set %v, but got %v", collectionName, len(page))
	}

}

func TestDeleteCollection(t *testing.T) {
	s, collection, _, _, ctx := InitializeCollectionTests(t)

	err := s.Images.DeleteCollection(ctx, collection)
	AssertNoError(t, err)

	images, _, err := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		im.OrderingArgs{},
		g.PaginationParams{},
		im.FetchMetaOnly)
	if len(images) != 0 {
		t.Fatalf("expected to retrieve no images, but got %v", len(images))
	}

	collection, err = s.Collections.Find(ctx, collection.Id)
	AssertErrorIs(t, err, e.ErrNotFound)

}
func TestDeleteCollectionWithAnnotations(t *testing.T) {
	s, collection, image, _, ctx := InitializeCollectionTests(t)
	label, _ := lbl.New("thelabel", "")

	s.Labels.Create(ctx, label)
	s.Images.Annotations.ApplyLabel(ctx, label, image)

	err := s.Images.DeleteCollection(ctx, collection)
	AssertNoError(t, err)

	images, _, err := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		im.OrderingArgs{},
		g.PaginationParams{},
		im.FetchMetaOnly)
	if len(images) != 0 {
		t.Fatalf("expected to retrieve no images, but got %v", len(images))
	}

	collection, err = s.Collections.Find(ctx, collection.Id)
	AssertErrorIs(t, err, e.ErrNotFound)

}

func TestDeleteImageFromCollection(t *testing.T) {

	s, collection, image, _, ctx := InitializeCollectionTests(t)

	err := s.Images.RemoveFromCollection(ctx, image)
	AssertNoError(t, err)

	images, _, _ := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		*im.NewImageDefaultOrderingArgs(),
		g.OneItemPaginationParams,
		im.FetchMetaOnly)
	if len(images) != 0 {
		t.Fatalf("expected to retrieve 0 images, but found %v", len(images))
	}

}

func TestAnnotationsShouldApplyToSpecifiedCollection(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)
	collectionName := "mycollection"
	labelName := "mylabel"
	collection, _ := clc.New(collectionName, "", "")

	notAnnotatedCollection, _ := clc.New("not-annotated-collection", "", "")
	image, _ := im.New(testPNGImage)
	label, _ := lbl.New(labelName, "")

	s.Collections.Create(ctx, collection)
	s.Collections.Create(ctx, notAnnotatedCollection)

	s.Images.Save(ctx, image, collection)
	s.Labels.Create(ctx, label)
	s.Images.Annotations.ApplyLabel(ctx, label, image)

	s.Images.AssignToCollection(ctx, image, notAnnotatedCollection)

	notAnnotatedImages, _, _ := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(notAnnotatedCollection.Id)),
		*im.NewImageDefaultOrderingArgs(),
		g.OneItemPaginationParams,
		im.FetchMetaOnly)
	if notAnnotatedImages[0].Annotations != nil {
		t.Fatal("expected to retrieve not annotated image")
	}

}

func InitImportImageTests(t *testing.T) (*a.App, context.Context, *im.Image, *clc.Collection, *clc.Collection, *im.BoundingBox) {
	s, _, ctx := a.NewTestApp(t, true)
	sourceCollectionName := "my-source-collection"
	destinationCollectionName := "my-destination-collection"
	labelName := "mylabel"

	sourceCollection, _ := clc.New(sourceCollectionName, "", "")
	destinationCollection, _ := clc.New(destinationCollectionName, "", "")
	image, _ := im.New(testPNGImage)
	label, _ := lbl.New(labelName, "")

	s.Collections.Create(ctx, sourceCollection)
	s.Collections.Create(ctx, destinationCollection)
	s.Labels.Create(ctx, label)

	s.Images.Save(ctx, image, sourceCollection)
	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)

	return &s, ctx, image, sourceCollection, destinationCollection, bbox

}

func TestImportImageInCollectionWithBoundingBoxes(t *testing.T) {

	s, ctx, srcImage, srcCollection, dstCollection, _ := InitImportImageTests(t)

	err := s.Images.ImportImage(ctx, srcImage, dstCollection.Id, im.ImportImageWithAnnotations)
	AssertNoError(t, err)
	retrievedDstImage, err := s.Images.Find(ctx, srcImage.Id, dstCollection.Id, im.FetchMetaOnly)
	if len(retrievedDstImage.BoundingBoxes) != 1 {
		t.Fatalf("expected to retrieve destination image with a bounding box, but got %v",
			len(retrievedDstImage.BoundingBoxes))
	}
	retrievedSrcImage, err := s.Images.Find(ctx, srcImage.Id, srcCollection.Id, im.FetchMetaOnly)
	if len(retrievedSrcImage.BoundingBoxes) != 1 {
		t.Fatalf("expected to retrieve source image with a bounding box, but got %v",
			len(retrievedSrcImage.BoundingBoxes))
	}
}

func TestImportImageInCollectionWithoutBoundingBoxes(t *testing.T) {

	s, ctx, srcImage, _, dstCollection, _ := InitImportImageTests(t)

	err := s.Images.ImportImage(ctx, srcImage, dstCollection.Id, im.ImportImageWithoutAnnotations)
	AssertNoError(t, err)
	retrievedDstImage, err := s.Images.Find(ctx, srcImage.Id, dstCollection.Id, im.FetchMetaOnly)
	AssertNoError(t, err)
	if len(retrievedDstImage.BoundingBoxes) != 0 {
		t.Fatalf("expected to retrieve image with no bounding box, but got %v",
			len(retrievedDstImage.BoundingBoxes))
	}
}

func TestImportDuplicateImageInCollectionShouldFail(t *testing.T) {

	s, ctx, srcImage, _, dstCollection, _ := InitImportImageTests(t)

	s.Images.ImportImage(ctx, srcImage, dstCollection.Id, im.ImportImageWithoutAnnotations)
	err := s.Images.ImportImage(ctx, srcImage, dstCollection.Id, im.ImportImageWithoutAnnotations)
	AssertErrorIs(t, err, e.ErrDuplication)
}

func TestPaginateCollections(t *testing.T) {

	tests := map[string]struct {
		nCollections         int
		page                 int64
		pageSize             int
		maxPageSize          int
		defaultPageSize      int
		expectedPageSize     int
		expectedTotalPages   int
		expectedTotalRecords int
	}{
		"one item page size 1": {nCollections: 1, page: 1, pageSize: 1, maxPageSize: 2, defaultPageSize: 1,
			expectedPageSize: 1, expectedTotalPages: 1, expectedTotalRecords: 1},
		"two items page size 1": {nCollections: 2, page: 1, pageSize: 1, maxPageSize: 2, defaultPageSize: 1,
			expectedPageSize: 1, expectedTotalPages: 2, expectedTotalRecords: 2},
		"two items page size 2": {nCollections: 2, page: 1, pageSize: 2, maxPageSize: 2, defaultPageSize: 1,
			expectedPageSize: 2, expectedTotalPages: 1, expectedTotalRecords: 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, _, ctx := a.NewTestApp(t, true)
			s.Collections.MaxPageSize = tc.maxPageSize

			for i := 0; i < tc.nCollections; i++ {
				clc, _ := clc.New(fmt.Sprintf("thename-%v", i+1), "", "")
				err := s.Collections.Create(ctx, clc)
				AssertNoError(t, err)
			}
			page, pageMeta, err := s.Collections.List(ctx,
				clc.AlphabeticalOrdering,
				g.PaginationParams{Page: tc.page, PageSize: tc.pageSize})
			AssertNoError(t, err)

			if len(page) != tc.expectedPageSize {
				t.Fatalf("expected to retrieve page of length 1, got %v", len(page))
			}

			nPages := int(pageMeta.TotalPages)
			nImages := pageMeta.TotalRecords
			AssertNoError(t, err)
			if nPages != tc.expectedTotalPages {
				t.Fatalf("expected pagination meta with total pages = %v, got %v.", tc.expectedTotalPages, nPages)
			}

			if int(nImages) != tc.expectedTotalRecords {
				t.Fatalf("expected pagination meta with total records = %v, got %v.", tc.expectedTotalRecords, nImages)
			}
		})
	}

}

func TestAddImageToCollectionShouldModifyUpdatedAtField(t *testing.T) {

	s, collection, _, clock, ctx := InitializeCollectionTests(t)
	currentCollection, _ := s.Collections.Find(ctx, collection.Id)
	lastTimestamp := currentCollection.UpdatedAt

	newImage, _ := im.New(testPNGImage)
	clock.Advance(1 * time.Hour)
	s.Images.Save(ctx, newImage, collection)

	retrievedCollection, _ := s.Collections.Find(ctx, collection.Id)
	if !retrievedCollection.UpdatedAt.After(lastTimestamp) {
		t.Fatalf("expected to modify updated_at field of collection. Last value was %v, got %v",
			lastTimestamp, retrievedCollection.UpdatedAt,
		)
	}

}

func TestAddAnnotationShouldModifyUpdatedAtField(t *testing.T) {

	s, collection, image, clock, ctx := InitializeCollectionTests(t)
	currentCollection, _ := s.Collections.Find(ctx, collection.Id)
	lastTimestamp := currentCollection.UpdatedAt

	clock.Advance(1 * time.Hour)

	label, _ := lbl.New("thelabel", "")
	s.Labels.Create(ctx, label)
	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)

	updatedCollection, _ := s.Collections.Find(ctx, collection.Id)
	if !updatedCollection.UpdatedAt.After(lastTimestamp) {
		t.Fatalf("expected updated timestamp, but got original %v and updated %v",
			updatedCollection.UpdatedAt, lastTimestamp)
	}

}

func TestModifyingAnnotationShouldModifyUpdatedAtField(t *testing.T) {

	s, collection, image, clock, ctx := InitializeCollectionTests(t)

	label, _ := lbl.New("thelabel", "")
	s.Labels.Create(ctx, label)
	bbox, _ := im.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)

	currentCollection, _ := s.Collections.Find(ctx, collection.Id)
	lastTimestamp := currentCollection.UpdatedAt

	bbox.Coords.Xc += 1

	clock.Advance(1 * time.Hour)
	s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)

	updatedCollection, _ := s.Collections.Find(ctx, collection.Id)
	if !updatedCollection.UpdatedAt.After(lastTimestamp) {
		t.Fatalf("expected updated timestamp, but got original %v and updated %v",
			updatedCollection.UpdatedAt, lastTimestamp)
	}

}

func TestCollectionsAreReturnedInAlphabeticalOrder(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "admin")
	nonOrderedNames := []string{"c", "b", "a"}
	for _, name := range nonOrderedNames {
		collection, _ := clc.New(name, "", "mygroup")
		err := s.Collections.Create(ctx, collection)
		AssertNoError(t, err)
	}

	orderedNames := []string{"a", "b", "c"}
	retrievedCollections, _, err := s.Collections.List(ctx, clc.AlphabeticalOrdering,
		g.PaginationParams{Page: 1, PageSize: 3})
	AssertNoError(t, err)
	for i, collection := range retrievedCollections {
		if collection.Name != orderedNames[i] {
			t.Fatalf("expected to retrieve site name %v, but got %v", orderedNames[i], collection.Name)
		}
	}

}
