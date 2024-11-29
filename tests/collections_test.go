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

	set := &m.Collection{Name: "myimageset"}
	set, err := s.Collections.Create(ctx, set)

	AssertNoError(t, err)
	AssertNoError(t, err)
	retrievedSet, err := s.Collections.GetOne(ctx, set.Id.String())
	AssertNoError(t, err)

	if retrievedSet.Name != set.Name {
		t.Fatalf("expected to retrieve identical names. Wanted %v, got %v",
			set.Name, retrievedSet.Name)
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
			set := &m.Collection{Name: tc.name}
			set, err := s.Collections.Create(ctx, set)
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
	setName := "myset"
	set, err := s.Collections.Create(ctx, &m.Collection{Name: setName})
	AssertNoError(t, err)
	image_in_set := &m.Image{Data: testImage}
	s.Images.Save(ctx, image_in_set)
	s.Collections.AppendImageToSet(ctx, image_in_set, set)

	s.Images.Save(ctx, &m.Image{Data: testImage})

	page, _, err := s.Images.GetPage(ctx, g.PaginationParams{Page: 1, PageSize: 4},
		&g.ImageFilterArgs{SetName: "myset"}, false)
	AssertNoError(t, err)

	if len(page) != 1 {
		t.Fatalf("expected to retrieve 1 image in set %v, but got %v", setName, len(page))
	}

	if page[0].Id != image_in_set.Id {
		t.Fatalf("expected to retrieve image appended in set %v, but got another one", setName)
	}

}
