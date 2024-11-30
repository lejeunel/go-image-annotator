package tests

import (
	"fmt"
	"github.com/go-test/deep"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	"testing"
)

func clearDateTime(image *m.Image) *m.Image {
	image.CreatedAt = ""
	image.UpdatedAt = ""
	return image
}

func TestSavingImageWithSHA256(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	testImageSHA256 := "cff295b60ef32bcd2e9a3c38eaf35dfdf78ffaf8bc95e655b682dd268329cfa1"
	image := &m.Image{Data: testImage, SHA256: testImageSHA256}
	collection := &m.Collection{Name: "thename"}
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	retrievedImage, _ := s.Images.Get(ctx, collection.Id.String(), image.Id.String(), true)

	diff := deep.Equal(image, retrievedImage)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical image structs, but got different fields: %v", diff))
	}

}

func TestSavingCorruptedImageWithSHA256ShouldFail(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	corruptSHA256 := "dff295b60ef32bcd2e9a3c38eaf35dfdf78ffaf8bc95e655b682dd268329cfa1"
	image := &m.Image{Data: testImage, SHA256: corruptSHA256}
	collection := &m.Collection{Name: "thename"}
	s.Collections.Create(ctx, collection)

	err := s.Images.Save(ctx, image, collection)
	AssertError(t, err)

}

func TestPaginateImages(t *testing.T) {

	tests := map[string]struct {
		nImages     int
		page        int64
		pageSize    int
		maxPageSize int
	}{
		"one image page size 1":  {nImages: 1, page: 1, pageSize: 1, maxPageSize: 2},
		"two images page size 1": {nImages: 2, page: 1, pageSize: 1, maxPageSize: 2},
		"two images page size 2": {nImages: 2, page: 1, pageSize: 2, maxPageSize: 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, ctx := NewTestApp(t, tc.maxPageSize)

			collection := &m.Collection{Name: "thename"}
			s.Collections.Create(ctx, collection)
			for i := 0; i < tc.nImages; i++ {
				image := &m.Image{Data: testImage}
				err := s.Images.Save(ctx, image, collection)
				AssertNoError(t, err)
			}
			page, pageMeta, err := s.Images.GetPage(ctx, g.PaginationParams{Page: tc.page, PageSize: tc.pageSize}, &g.ImageFilterArgs{}, false)
			AssertNoError(t, err)

			if len(page) != tc.pageSize {
				t.Fatalf("expected to retrieve page of length 1, got %v", len(page))

			}

			if pageMeta.TotalPages != tc.nImages/tc.pageSize {
				t.Fatalf("expected pagination meta with total pages = %v, got %v", tc.pageSize, pageMeta.TotalPages)
			}

			if int(pageMeta.TotalRecords) != tc.nImages {
				t.Fatalf("expected pagination meta with total records = 2, got %v", pageMeta.TotalRecords)
			}
		})
	}

}
