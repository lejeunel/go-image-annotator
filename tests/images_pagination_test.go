package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	g "datahub/generic"
	"testing"
)

type ImagePaginationTestEnv struct {
	app        *a.App
	collection *clc.Collection
	ctx        context.Context
}

func NewImagePaginationTestEnv(t *testing.T, maxPageSize int, nImages int) *ImagePaginationTestEnv {
	s, _, ctx := a.NewTestApp(t, true)
	s.Images.MaxPageSize = maxPageSize

	collection, _ := clc.New("thename", "", "")
	s.Collections.Create(ctx, collection)
	for range nImages {
		image, _ := im.New(testPNGImage)
		err := s.Images.Save(ctx, image, collection)
		AssertNoError(t, err)
	}

	return &ImagePaginationTestEnv{&s, collection, ctx}

}

type PaginateImagesTestCase struct {
	testName             string
	nImages              int
	page                 int64
	pageSize             int
	maxPageSize          int
	defaultPageSize      int
	expectedPageSize     int
	expectedTotalPages   int
	expectedTotalRecords int
}

func assertImagePagination(t *testing.T, page []im.Image, meta *g.PaginationMeta, tc PaginateImagesTestCase) {
	if len(page) != tc.expectedPageSize {
		t.Fatalf("expected to retrieve page of length 1, got %v", len(page))
	}

	nPages := int(meta.TotalPages)
	nImages := meta.TotalRecords
	if nPages != tc.expectedTotalPages {
		t.Fatalf("expected pagination meta with total pages = %v, got %v", tc.expectedTotalPages, nPages)
	}

	if int(nImages) != tc.expectedTotalRecords {
		t.Fatalf("expected pagination meta with total records = %v, got %v", tc.expectedTotalRecords, nImages)
	}

}

func TestPaginateImages(t *testing.T) {

	tests := []PaginateImagesTestCase{
		{"one item page size 1", 1, 1, 1, 2, 1, 1, 1, 1},
		{"two items page size 1", 2, 1, 1, 2, 1, 1, 2, 2},
		{"two items page size 2", 2, 1, 2, 2, 1, 2, 1, 2},
	}

	for _, tc := range tests {
		env := NewImagePaginationTestEnv(t, tc.maxPageSize, tc.nImages)
		page, meta, err := env.app.Images.List(env.ctx,
			*im.NewImageFilter(im.WithCollectionId(env.collection.Id)),
			im.OrderingArgs{},
			g.PaginationParams{Page: tc.page, PageSize: tc.pageSize},
			im.FetchMetaOnly)
		AssertNoError(t, err)
		assertImagePagination(t, page, meta, tc)

	}

}
