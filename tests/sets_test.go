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

func TestCreateSet(t *testing.T) {
	s, ctx := NewTestComponents(t, 2)

	set := &m.Set{Name: "myimageset"}
	set, err := s.Sets.Create(ctx, set)

	AssertNoError(t, err)
	AssertNoError(t, err)
	retrievedSet, err := s.Sets.GetOne(ctx, set.Id.String())
	AssertNoError(t, err)

	if retrievedSet.Name != set.Name {
		t.Fatalf("expected to retrieve identical names. Wanted %v, got %v",
			set.Name, retrievedSet.Name)
	}

}

func TestRetrieveImagesOfSet(t *testing.T) {

	s, ctx := NewTestComponents(t, 2)
	setName := "myset"
	set, err := s.Sets.Create(ctx, &m.Set{Name: setName})
	AssertNoError(t, err)
	image_in_set, _ := s.Images.Save(ctx, &m.Image{Data: testImage})
	s.Sets.AppendImageToSet(ctx, image_in_set, set)

	s.Images.Save(ctx, &m.Image{Data: testImage})

	page, _, err := s.Images.GetPage(ctx, g.PaginationParams{Page: 1, PageSize: 4}, &g.ImageFilterArgs{SetName: "myset"}, false)
	AssertNoError(t, err)

	if len(page) != 1 {
		t.Fatalf("expected to retrieve 1 image in set %v, but got %v", setName, len(page))
	}

	if page[0].Id != image_in_set.Id {
		t.Fatalf("expected to retrieve image appended in set %v, but got another one", setName)
	}

}
