package tests

import (
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	"testing"
)

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

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

	page, _, err := s.Images.GetPage(ctx, g.PaginationParams{Page: 1, PageSize: 4},
		&g.ImageFilterArgs{CollectionName: "myset"}, false)
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

	images, _, _ := s.Images.GetPage(ctx, g.PaginationParams{}, &g.ImageFilterArgs{}, false)
	if len(images) > 0 {
		t.Fatal("expected to retrieve 0 images, but found some")
	}

}

func TestCloningCollectionsShouldNotDuplicateImages(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	collectionName := "mycollection"
	collection := &m.Collection{Name: collectionName}
	s.Collections.Create(ctx, collection)
	image := &m.Image{Data: testImage}
	s.Images.Save(ctx, image, collection)

	s.Collections.Clone(ctx, collection, &m.Collection{Name: "theclone"})
	imagesOfClone, _, _ := s.Images.GetPage(ctx, g.PaginationParams{},
		&g.ImageFilterArgs{CollectionName: "theclone"}, false)

	images, _, _ := s.Images.GetPage(ctx, g.PaginationParams{}, &g.ImageFilterArgs{}, false)
	if len(images) != 1 {
		t.Fatalf("expected to get 1 image accross all collections, but found %v", len(images))
	}

	if len(imagesOfClone) != 1 {
		t.Fatalf("expected to retrieve 1 image in cloned collection, but found %v", len(imagesOfClone))
	}

	if imagesOfClone[0].Id != images[0].Id {
		t.Fatalf("expected to retrieve images in cloned collection with identical id, but it is different: %v", imagesOfClone[0].Id)
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
	imagesInMerged, _, _ := s.Images.GetPage(ctx, g.PaginationParams{Page: 1, PageSize: 4},
		&g.ImageFilterArgs{CollectionName: "first"}, false)

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
	imagesInMerged, _, _ := s.Images.GetPage(ctx, g.PaginationParams{Page: 1, PageSize: 4},
		&g.ImageFilterArgs{CollectionName: "new-collection"}, false)

	if len(imagesInMerged) != 1 {
		t.Fatalf("expected to retrieve 1 image in merged collection, got %v",
			len(imagesInMerged))
	}

}
