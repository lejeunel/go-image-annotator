package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	loc "datahub/domain/locations"
	e "datahub/errors"
	g "datahub/generic"
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

func InitPatchTests(t *testing.T) (a.App, *im.BaseImage, *loc.Site, *loc.Camera, *clk.FakeClock, context.Context) {
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
	s.Images.AssignCamera(ctx, camera.Id, image.Id)

	base, _ := s.Images.GetBase(ctx, image.Id, im.FetchMetaOnly)

	return s, base, site, camera, clock, ctx

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
	for range nImages {
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

type ImageTypeTestCase struct {
	Type    string
	WantErr bool
}

func TestSavingImageWithType(t *testing.T) {

	tests := []ImageTypeTestCase{{"thermal", false}, {"rgb", false},
		{"gray", false}, {"undefined-type", true}}

	for _, tc := range tests {
		s, _, ctx := a.NewTestApp(t, true)

		image, _ := im.New(testJPGImage)
		image.Type = tc.Type
		collection, _ := clc.New("thename", "", "")
		s.Collections.Create(ctx, collection)
		err := s.Images.Save(ctx, image, collection)
		if tc.WantErr {
			if err == nil {
				t.Fatalf("expected an error with image of type %v, but got none", tc.Type)
			}
			return
		}
		if err != nil {
			t.Fatalf("unexpected error with image type %v: %v", tc.Type, err)
		}

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
