package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	loc "datahub/domain/locations"
	e "datahub/errors"
	g "datahub/generic"
	"fmt"
	clk "github.com/jonboulle/clockwork"
	"testing"
	"time"
)

func TestCreateImage(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	image, err := im.New(testJPGImage)
	AssertNoError(t, err)
	collection, _ := clc.New("mycollection", "mydescription", "mygroup")
	s.Collections.Create(ctx, collection)
	err = s.Images.Save(ctx, image, collection)
	AssertNoError(t, err)

	res, err := s.Images.Find(ctx, image.Id, collection.Id, im.FetchWithRawData)
	AssertNoError(t, err)
	AssertDeepEqual(t, image, res, "image")
}

func TestFetchingNonExistingImageShouldFail(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("mycollection", "mydescription", "mygroup")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	wrongId := im.NewImageId()
	_, err := s.Images.Find(ctx, *wrongId, collection.Id, im.FetchWithRawData)
	AssertErrorIs(t, err, e.ErrNotFound)
}
func TestFetchingImageThatExistsInTwoCollections(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	image, _ := im.New(testJPGImage)
	firstCollection, _ := clc.New("mycollection", "mydescription", "mygroup")
	secondCollection, _ := clc.New("mysecondcollection", "mydescription", "mygroup")
	s.Collections.Create(ctx, firstCollection)
	s.Collections.Create(ctx, secondCollection)
	s.Images.Save(ctx, image, firstCollection)
	s.Images.Save(ctx, image, secondCollection)

	retrieved, _ := s.Images.Find(ctx, image.Id, secondCollection.Id, im.FetchMetaOnly)
	if retrieved.Collection.Name != "mysecondcollection" {
		t.Fatalf("expected to retrieve image in collection %v but got %v",
			"mysecondcollection", retrieved.Collection.Name)
	}
}

func TestUpdateImage(t *testing.T) {

	timeLayout := "2006-01-02T15:04:05.000Z"
	tests := []struct {
		testName   string
		capturedAt string
		site       string
		camera     string
		type_      string
		success    bool
	}{
		{testName: "non-existing site", site: "non-existing-site", camera: "my-camera",
			success: false, capturedAt: "2006-01-02T15:04:05.000Z", type_: "rgb"},
		{testName: "non-existing camera", site: "my-site", camera: "non-existing-camera",
			success: false, capturedAt: "2006-01-02T15:04:05.000Z", type_: "rgb"},
		{testName: "invalid timestamp", site: "my-site", camera: "my-camera",
			success: false, capturedAt: "2006-01-02T15:04:00Z", type_: "rgb"},
		{testName: "invalid type", site: "my-site", camera: "my-camera",
			success: false, capturedAt: "2006-01-02T15:04:00Z", type_: "asdf"},
		{testName: "changing type", site: "my-site", camera: "my-camera",
			success: true, capturedAt: "2006-01-02T15:04:05.000Z", type_: "thermal"},
		{testName: "changing camera", site: "my-other-site", camera: "my-other-camera",
			success: true, capturedAt: "2006-01-02T15:04:05.000Z", type_: "thermal"},
	}

	for _, tc := range tests {
		s, _, ctx := a.NewTestApp(t, true)

		image, _ := im.New(testJPGImage)
		collection, _ := clc.New("my-collection", "", "")
		s.Collections.Create(ctx, collection)
		s.Images.Save(ctx, image, collection)

		site, _ := loc.NewSite("my-site", "thegroup")
		camera, _ := loc.NewCamera("my-camera", site, "")
		secondCamera, _ := loc.NewCamera("my-second-camera", site, "")
		s.Locations.SaveSite(ctx, site)
		s.Locations.SaveCamera(ctx, camera)
		s.Locations.SaveCamera(ctx, secondCamera)

		otherSite, _ := loc.NewSite("my-other-site", "thegroup")
		otherCamera, _ := loc.NewCamera("my-other-camera", otherSite, "")
		s.Locations.SaveSite(ctx, otherSite)
		s.Locations.SaveCamera(ctx, otherCamera)

		image, err := s.Images.Update(ctx, image.Id, im.ImageUpdatables{Site: tc.site, Camera: tc.camera, CapturedAt: tc.capturedAt,
			Type_: tc.type_})
		if tc.success {
			if err != nil {
				t.Error(fmt.Printf("%s: expected no error but got %v", tc.testName, err))
			}
			wantT, _ := time.Parse(timeLayout, tc.capturedAt)
			if image.CapturedAt != wantT {
				t.Fatalf("%v: expected captured_at to be %v, but got %v", tc.testName, wantT, image.CapturedAt)
			}
			if image.GetSiteName() != tc.site {
				t.Fatalf("expected site to be %v, but got %v", tc.site, image.GetSiteName())
			}
			if image.GetCameraName() != tc.camera {
				t.Fatalf("expected camera to be %v, but got %v", tc.camera, image.GetCameraName())
			}
			if image.Type != tc.type_ {
				t.Fatalf("expected type to be %v, but got %v", tc.type_, image.Type)
			}

		} else {
			if err == nil {
				t.Error(fmt.Printf("%s: expected error but got no error", tc.testName))
			}
		}

	}

}

func InitPatchTests(t *testing.T) (a.App, *im.Image, *loc.Site, *loc.Camera, *clk.FakeClock, context.Context) {
	s, clock, ctx := a.NewTestApp(t, true)

	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("my-collection", "", "")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	site, _ := loc.NewSite("my-site", "thegroup")
	camera, _ := loc.NewCamera("my-camera", site, "")
	secondCamera, _ := loc.NewCamera("my-second-camera", site, "")
	s.Locations.SaveSite(ctx, site)
	s.Locations.SaveCamera(ctx, camera)
	s.Locations.SaveCamera(ctx, secondCamera)
	s.Images.AssignCamera(ctx, camera, image)

	return s, image, site, camera, clock, ctx

}

func TestPatchCameraOnImage(t *testing.T) {
	s, image, _, _, _, ctx := InitPatchTests(t)
	patch := []g.JSONPatch{{Operation: "replace", Path: "/camera", Value: "my-second-camera"}}
	image, err := s.Images.Patch(ctx, image.Id, patch)
	AssertNoError(t, err)
	if image.GetCameraName() != "my-second-camera" {
		t.Fatalf("expected to change camera name, but did not. Got %v instead", image.GetCameraName())
	}

}

func TestOrderingImagesByCapturedAt(t *testing.T) {
	s, clock, ctx := a.NewTestApp(t, true)
	nImages := 3
	collection, _ := clc.New("thename", "", "mygroup")
	s.Collections.Create(ctx, collection)
	for i := 0; i < nImages; i++ {
		image, _ := im.New(testPNGImage)
		image.CapturedAt = clock.Now()
		s.Images.Save(ctx, image, collection)
		clock.Advance(2 * time.Minute)
	}
	ascendingImages, _, _ := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		*im.NewAscendingImageCapturedOrder(),
		g.PaginationParams{Page: 1, PageSize: 2},
		im.FetchMetaOnly)
	descendingImages, _, _ := s.Images.List(ctx,
		*im.NewImageFilter(im.WithCollectionId(collection.Id)),
		*im.NewDescendingImageCapturedOrder(),
		g.PaginationParams{Page: 1, PageSize: 2},
		im.FetchMetaOnly)

	if ascendingImages[0].CapturedAt.Before(ascendingImages[1].CapturedAt) == false {
		t.Fatalf("expected to retrieve images with captured_at field in ascending order, got %v and %v",
			ascendingImages[0].CapturedAt, ascendingImages[1].CapturedAt)
	}

	if descendingImages[0].CapturedAt.Before(descendingImages[1].CapturedAt) == true {
		t.Fatalf("expected to retrieve images with captured_at field in descending order, got %v and %v",
			descendingImages[0].CapturedAt, descendingImages[1].CapturedAt)
	}

}

func TestFilterImagesByCollectionName(t *testing.T) {

	s, _, ctx := a.NewTestApp(t, true)
	firstCollection, _ := clc.New("first-collection", "", "")
	secondCollection, _ := clc.New("second-collection", "", "")
	s.Collections.Create(ctx, firstCollection)
	s.Collections.Create(ctx, secondCollection)

	firstImage, _ := im.New(testPNGImage)
	secondImage, _ := im.New(testJPGImage)
	s.Images.Save(ctx, firstImage, firstCollection)
	s.Images.Save(ctx, secondImage, secondCollection)

	retrievedImages, paginationMeta, err := s.Images.List(ctx, *im.NewImageFilter(im.WithCollectionName("first-collection")),
		im.OrderingArgs{}, g.PaginationParams{Page: 1, PageSize: 2}, im.FetchMetaOnly)
	AssertNoError(t, err)
	if paginationMeta.TotalRecords != 1 {
		t.Fatalf("expected to get total records of 1, got %v", paginationMeta.TotalRecords)
	}
	if retrievedImages[0].Id != firstImage.Id {
		t.Fatal("expected to retrieve first image")
	}

}

func TestPaginateImages(t *testing.T) {

	tests := map[string]struct {
		nImages              int
		page                 int64
		pageSize             int
		maxPageSize          int
		defaultPageSize      int
		expectedPageSize     int
		expectedTotalPages   int
		expectedTotalRecords int
	}{
		"one item page size 1": {nImages: 1, page: 1, pageSize: 1, maxPageSize: 2, defaultPageSize: 1,
			expectedPageSize: 1, expectedTotalPages: 1, expectedTotalRecords: 1},
		"two items page size 1": {nImages: 2, page: 1, pageSize: 1, maxPageSize: 2, defaultPageSize: 1,
			expectedPageSize: 1, expectedTotalPages: 2, expectedTotalRecords: 2},
		"two items page size 2": {nImages: 2, page: 1, pageSize: 2, maxPageSize: 2, defaultPageSize: 1,
			expectedPageSize: 2, expectedTotalPages: 1, expectedTotalRecords: 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, _, ctx := a.NewTestApp(t, true)
			s.Images.MaxPageSize = tc.maxPageSize

			collection := &clc.Collection{Name: "thename"}
			s.Collections.Create(ctx, collection)
			for i := 0; i < tc.nImages; i++ {
				image, _ := im.New(testPNGImage)
				err := s.Images.Save(ctx, image, collection)
				AssertNoError(t, err)
			}
			page, pageMeta, err := s.Images.List(ctx,
				*im.NewImageFilter(im.WithCollectionId(collection.Id)),
				im.OrderingArgs{},
				g.PaginationParams{Page: tc.page, PageSize: tc.pageSize},
				im.FetchMetaOnly)
			AssertNoError(t, err)

			if len(page) != tc.expectedPageSize {
				t.Fatalf("expected to retrieve page of length 1, got %v", len(page))
			}

			nPages := int(pageMeta.TotalPages)
			nImages := pageMeta.TotalRecords
			AssertNoError(t, err)
			if nPages != tc.expectedTotalPages {
				t.Fatalf("expected pagination meta with total pages = %v, got %v", tc.expectedTotalPages, nPages)
			}

			if int(nImages) != tc.expectedTotalRecords {
				t.Fatalf("expected pagination meta with total records = %v, got %v", tc.expectedTotalRecords, nImages)
			}
		})
	}

}

func TestSavingImageWithType(t *testing.T) {

	tests := map[string]struct {
		Type    string
		IsValid bool
	}{
		"thermal":        {Type: "thermal", IsValid: true},
		"rgb":            {Type: "rgb", IsValid: true},
		"gray":           {Type: "gray", IsValid: true},
		"undefined type": {Type: "", IsValid: true},
		"whatever":       {Type: "whatever", IsValid: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, _, ctx := a.NewTestApp(t, true)

			image, _ := im.New(testJPGImage)
			image.Type = tc.Type
			collection, _ := clc.New("thename", "", "")
			s.Collections.Create(ctx, collection)
			err := s.Images.Save(ctx, image, collection)
			if tc.IsValid {
				if err != nil {
					t.Fatalf("expected to get no error with image of type %v, but got one", tc.Type)
				}
			} else {
				if err == nil {
					t.Fatalf("expected to get an error with image of type %v, but got none", tc.Type)
				}

			}
		})
	}

}

func TestRawImageURL(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("mycollection", "mydescription", "mygroup")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	URLBuilder := im.NewRawImageURLBuilder(
		g.NewAPIURLBuilder("v1"),
		"raw-image")
	response := im.NewImageResponse(image, URLBuilder)

	want := URLBuilder.Build(image.Id)
	got := response.RawURL
	if got != want {
		t.Fatalf("expected to retrieve raw image url %v, but got %v",
			want, got,
		)
	}

}
