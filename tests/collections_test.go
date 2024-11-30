package tests

import (
	"fmt"
	"github.com/go-test/deep"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	"testing"
)

func TestCreateCollection(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	clc := &m.Collection{Name: "myimageset"}
	err := s.Collections.Create(ctx, clc)

	AssertNoError(t, err)
	retrievedSet, err := s.Collections.Get(ctx, clc.Id.String())
	AssertNoError(t, err)

	if retrievedSet.Name != clc.Name {
		t.Fatalf("expected to retrieve identical names. Wanted %v, got %v",
			clc.Name, retrievedSet.Name)
	}
}

func TestCollectionDuplicateNameShouldFail(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	clc := &m.Collection{Name: "myimageset"}
	newClc := &m.Collection{Name: "myimageset"}
	s.Collections.Create(ctx, clc)
	err := s.Collections.Create(ctx, newClc)

	AssertError(t, err)
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

		s, ctx := NewTestApp(t, 2)
		t.Run(name, func(t *testing.T) {
			collection := &m.Collection{Name: tc.name}
			err := s.Collections.Create(ctx, collection)
			if tc.isValid {
				AssertNoError(t, err)
			} else {
				AssertError(t, err)
			}
		})
	}

}

func TestRetrieveImagesOfCollection(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "myset"
	collection := &m.Collection{Name: collectionName}
	err := s.Collections.Create(ctx, collection)
	AssertNoError(t, err)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, collection)

	page, _, err := s.Images.GetPage(ctx, collection.Id.String(), g.PaginationParams{Page: 1, PageSize: 4},
		false)
	AssertNoError(t, err)

	if len(page) != 1 {
		t.Fatalf("expected to retrieve 1 image in set %v, but got %v", collectionName, len(page))
	}

}

func TestOrphanImagesShouldBeDeleted(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, collection)

	s.Collections.Delete(ctx, collection)

	images, _, _ := s.Images.GetPage(ctx, collection.Id.String(), g.PaginationParams{}, false)
	if len(images) > 0 {
		t.Fatal("expected to retrieve 0 images, but found some")
	}

}

func TestAnnotationsShouldApplyToSpecifiedCollection(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	labelName := "mylabel"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: labelName}
	s.Images.Save(ctx, image, collection)
	s.Annotations.CreateLabel(ctx, label)
	s.Annotations.ApplyLabelToImage(ctx, label, image, collection)

	notAnnotatedCollection := &m.Collection{Name: "not-annotated-collection"}
	s.Collections.Create(ctx, notAnnotatedCollection)
	s.Collections.CollectionRepo.AssignImageToCollection(ctx, image, notAnnotatedCollection)

	notAnnotatedImages, _, _ := s.Images.GetPage(ctx, notAnnotatedCollection.Id.String(), g.PaginationParams{}, false)
	if notAnnotatedImages[0].Annotations != nil {
		t.Fatal("expected to retrieve not annotated image")
	}

}

func TestCloningCollectionsShouldAlsoCloneAnnotations(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	cloneName := "theclone"
	labelName := "mylabel"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: labelName}
	s.Images.Save(ctx, image, collection)
	s.Annotations.CreateLabel(ctx, label)
	s.Annotations.ApplyLabelToImage(ctx, label, image, collection)
	bbox := m.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)

	s.Annotations.ApplyBoundingBoxToImage(ctx, bbox, image, collection)

	clone := &m.Collection{Name: cloneName}
	s.Collections.Clone(ctx, collection, clone)
	cloneImages, _, _ := s.Images.GetPage(ctx, clone.Id.String(), g.PaginationParams{},
		false)

	diff := deep.Equal(image.Annotations, cloneImages[0].Annotations)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical image annotations, but got different fields: %v", diff))
	}

	diff = deep.Equal(image.BoundingBoxes, cloneImages[0].BoundingBoxes)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical image bounding boxes, but got different fields: %v", diff))
	}

}

func TestCloningCollectionsShouldNotDuplicateImages(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, collection)

	clonedCollection := &m.Collection{Name: "theclone"}

	s.Collections.Clone(ctx, collection, clonedCollection)
	imagesOfClone, _, _ := s.Images.GetPage(ctx, clonedCollection.Id.String(),
		g.PaginationParams{}, false)
	imagesOrigin, _, _ := s.Images.GetPage(ctx, collection.Id.String(),
		g.PaginationParams{}, false)

	if len(imagesOfClone) != 1 {
		t.Fatalf("expected to retrieve 1 image in cloned collection, but found %v", len(imagesOfClone))
	}

	if imagesOfClone[0].Id != imagesOrigin[0].Id {
		t.Fatalf("expected to retrieve images in cloned collection with identical id, but it is different: %v VS %v",
			imagesOfClone[0].Id, imagesOrigin[0].Id)
	}

}

func TestMergingCollections(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	firstCollection := &m.Collection{Name: "first"}
	secondCollection := &m.Collection{Name: "second"}
	s.Collections.Create(ctx, firstCollection)
	s.Collections.Create(ctx, secondCollection)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, firstCollection)
	s.Images.Save(ctx, image, secondCollection)

	s.Collections.Merge(ctx, secondCollection, firstCollection)
	imagesInMerged, _, _ := s.Images.GetPage(ctx, firstCollection.Id.String(),
		g.PaginationParams{Page: 1, PageSize: 4}, false)

	if len(imagesInMerged) != 2 {
		t.Fatalf("expected to retrieve 2 images in merging destination collection, got %v",
			len(imagesInMerged))
	}
}

func TestMergingCollectionsShouldSkipDuplicateImages(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	commonImage := &m.Image{Data: testImage}
	s.Images.Save(ctx, commonImage, collection)

	newCollection := &m.Collection{Name: "new-collection"}
	s.Collections.Create(ctx, newCollection)
	s.Collections.CollectionRepo.AssignImageToCollection(ctx, commonImage, newCollection)

	s.Collections.Merge(ctx, collection, newCollection)
	imagesInMerged, _, _ := s.Images.GetPage(ctx, newCollection.Id.String(),
		g.PaginationParams{Page: 1, PageSize: 4},
		false)

	if len(imagesInMerged) != 1 {
		t.Fatalf("expected to retrieve 1 image in merged collection, got %v",
			len(imagesInMerged))
	}

}
