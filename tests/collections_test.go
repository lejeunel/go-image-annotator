package tests

import (
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

func TestDeleteOrphanImagesShouldBeDeletedAsWell(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, collection)

	s.Collections.Delete(ctx, collection)

	images, _ := s.Images.GetOrdinal(ctx, collection.Id.String(), 0, false)
	if images != nil {
		t.Fatal("expected to retrieve 0 images, but found some")
	}

}

func TestDeleteImageFromCollection(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, collection)

	err := s.Collections.RemoveImage(ctx, image, collection)
	AssertNoError(t, err)

	images, _ := s.Images.GetOrdinal(ctx, collection.Id.String(), 0, false)
	if images != nil {
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

	notAnnotatedImage, _ := s.Images.GetOrdinal(ctx, notAnnotatedCollection.Id.String(), 0, false)
	if notAnnotatedImage.Annotations != nil {
		t.Fatal("expected to retrieve not annotated image")
	}

}

func TestDeepCloneShouldAlsoCloneAnnotations(t *testing.T) {

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
	s.Collections.Clone(ctx, collection, clone, true)
	cloneImage, _ := s.Images.GetOrdinal(ctx, clone.Id.String(), 0,
		false)

	if cloneImage.Annotations == nil {
		t.Fatal("expected to retrieve image annotations, but got none")
	}

	if cloneImage.BoundingBoxes == nil {
		t.Fatal("expected to retrieve image with bounding boxes, but got none")
	}

}

func TestCloneShouldSkipDuplicateImages(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, collection)

	clonedCollection := &m.Collection{Name: "theclone"}

	s.Collections.Clone(ctx, collection, clonedCollection, false)
	imageOfClone, _ := s.Images.GetOrdinal(ctx, clonedCollection.Id.String(), 0, false)
	imageOrigin, _ := s.Images.GetOrdinal(ctx, collection.Id.String(), 0, false)

	if imageOfClone == nil {
		t.Fatal("expected to retrieve 1 image in cloned collection, but found none")
	}

	if imageOfClone.Id != imageOrigin.Id {
		t.Fatalf("expected to retrieve images in cloned collection with identical id, but it is different: %v VS %v",
			imageOfClone.Id, imageOrigin.Id)
	}

}

func TestMergeCollections(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	firstCollection := &m.Collection{Name: "first"}
	secondCollection := &m.Collection{Name: "second"}
	s.Collections.Create(ctx, firstCollection)
	s.Collections.Create(ctx, secondCollection)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, firstCollection)
	s.Images.Save(ctx, image, secondCollection)

	s.Collections.Merge(ctx, secondCollection, firstCollection, false)
	imagesInMerged, _, _ := s.Images.GetPage(ctx, firstCollection.Id.String(),
		g.PaginationParams{Page: 1, PageSize: 4}, false)

	if len(imagesInMerged) != 2 {
		t.Fatalf("expected to retrieve 2 images in merging destination collection, got %v",
			len(imagesInMerged))
	}
}

func TestMergeShouldSkipDuplicateImages(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	commonImage := &m.Image{Data: testImage}
	s.Images.Save(ctx, commonImage, collection)

	newCollection := &m.Collection{Name: "new-collection"}
	s.Collections.Create(ctx, newCollection)
	s.Collections.CollectionRepo.AssignImageToCollection(ctx, commonImage, newCollection)

	s.Collections.Merge(ctx, collection, newCollection, false)
	imagesInMerged, _, _ := s.Images.GetPage(ctx, newCollection.Id.String(),
		g.PaginationParams{Page: 1, PageSize: 4},
		false)

	if len(imagesInMerged) != 1 {
		t.Fatalf("expected to retrieve 1 image in merged collection, got %v",
			len(imagesInMerged))
	}

}

func TestDeepMergeShouldAlsoCopyAnnotations(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	labelName := "mylabel"

	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	collectionToMerge := &m.Collection{Name: "mycollectiontomerge"}
	s.Collections.Create(ctx, collectionToMerge)

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: labelName}
	s.Images.Save(ctx, image, collectionToMerge)
	s.Annotations.CreateLabel(ctx, label)
	s.Annotations.ApplyLabelToImage(ctx, label, image, collectionToMerge)
	bbox := m.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)
	s.Annotations.ApplyBoundingBoxToImage(ctx, bbox, image, collectionToMerge)

	s.Collections.Merge(ctx, collectionToMerge, collection, true)
	destinationImages, _, _ := s.Images.GetPage(ctx, collection.Id.String(), g.PaginationParams{},
		false)

	if destinationImages[0].Annotations == nil {
		t.Fatal("expected to retrieve image annotations, but got none")
	}

	if destinationImages[0].BoundingBoxes == nil {
		t.Fatal("expected to retrieve image with bounding boxes, but got none")
	}

}
