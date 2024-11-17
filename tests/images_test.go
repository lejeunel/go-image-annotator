package tests

import (
	"bytes"
	m "go-image-annotator/models"
	"testing"
)

func TestSaveAndRetrieveImage(t *testing.T) {
	s, ctx := NewTestComponents(t)

	image := &m.Image{Data: testImage}

	image, err := s.Images.Save(ctx, image)
	AssertNoError(t, err)

	retrievedImage, err := s.Images.GetOne(ctx, image.Id.String())

	if image.Width != retrievedImage.Width {
		t.Fatalf("expected to retrieve identical widths. Wanted %v, got %v",
			image.Width, retrievedImage.Width)
	}

	if !bytes.Equal(testImage, retrievedImage.Data) {
		t.Fatalf("expected to retrieve identical data, but they are different. Wanted %v, got %v", testImage, retrievedImage.Data)
	}

}

func TestSavingImageWithSHA256(t *testing.T) {
	s, ctx := NewTestComponents(t)

	testImageSHA256 := "cff295b60ef32bcd2e9a3c38eaf35dfdf78ffaf8bc95e655b682dd268329cfa1"
	image := &m.Image{Data: testImage, SHA256: testImageSHA256}

	_, err := s.Images.Save(ctx, image)
	AssertNoError(t, err)

}

func TestSavingCorruptedImageWithSHA256ShouldFail(t *testing.T) {
	s, ctx := NewTestComponents(t)

	corruptSHA256 := "dff295b60ef32bcd2e9a3c38eaf35dfdf78ffaf8bc95e655b682dd268329cfa1"
	image := &m.Image{Data: testImage, SHA256: corruptSHA256}

	_, err := s.Images.Save(ctx, image)
	AssertError(t, err)

}
