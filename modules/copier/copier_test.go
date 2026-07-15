package copier

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingSourceCollectionShouldFail(t *testing.T) {
	cp := NewTestingCopier()
	cp.CollectionRepo = &FakeCollectionRepo{MissingCollection: true}
	_, err := cp.Copy(Request{})
	assert.Error(t, err)
}
func TestHandleImageStoreFail(t *testing.T) {
	cp := NewTestingCopier()
	cp.Store = &FakeImageStore{Err: e.ErrNotFound}
	_, err := cp.Copy(Request{})
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func PrepareCopying() (Copier, clc.Collection, clc.Collection, im.Image, lbl.Label) {
	cp := NewTestingCopier()
	srcCollection := clc.NewCollection(clc.NewCollectionId(), "src-collection")
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection")
	srcImage := im.NewImage(im.NewImageId(), srcCollection)
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")

	return cp, srcCollection, dstCollection, srcImage, label

}

func TestImageLabelsNotCopiedWhenShallow(t *testing.T) {
	cp, dst, src, image, label := PrepareCopying()
	image.AddLabel(label)
	cp.Store = &FakeImageStore{Return: &image}
	repo := &FakeAnnotationRepo{}
	cp.AnnotationRepo = repo
	_, err := cp.Copy(Request{ImageId: image.Id, SourceCollection: src.Name,
		DestinationCollection: dst.Name, Deep: false})
	assert.NoError(t, err)
	assert.Equal(t, 0, repo.NumLabelsAdded)
}

func TestBoundingBoxesAreCopied(t *testing.T) {
	cp, dst, src, image, label := PrepareCopying()
	image.AddBoundingBox(a.NewBoundingBox(a.NewAnnotationId(), 0., 0., 1., 2., label))
	cp.Store = &FakeImageStore{Return: &image}
	repo := &FakeAnnotationRepo{}
	cp.AnnotationRepo = repo
	_, err := cp.Copy(Request{ImageId: image.Id, SourceCollection: src.Name,
		DestinationCollection: dst.Name, Deep: true})
	assert.NoError(t, err)
	assert.Equal(t, 1, repo.NumBoundingboxesAdded)
}

func TestImageLabelsAreCopied(t *testing.T) {
	cp, dst, src, image, label := PrepareCopying()
	image.AddLabel(label)
	cp.Store = &FakeImageStore{Return: &image}
	repo := &FakeAnnotationRepo{}
	cp.AnnotationRepo = repo
	_, err := cp.Copy(Request{ImageId: image.Id, SourceCollection: src.Name,
		DestinationCollection: dst.Name, Deep: true})
	assert.NoError(t, err)
	assert.Equal(t, 1, repo.NumLabelsAdded)
}

func TestPolygonsAreCopied(t *testing.T) {
	cp, dst, src, image, label := PrepareCopying()
	image.AddPolygon(a.NewPolygon(a.NewAnnotationId(),
		a.Points{Coordinates: [][2]float32{{0, 0}, {1, 1}}},
		label))
	cp.Store = &FakeImageStore{Return: &image}
	repo := &FakeAnnotationRepo{}
	cp.AnnotationRepo = repo
	_, err := cp.Copy(Request{ImageId: image.Id, SourceCollection: src.Name,
		DestinationCollection: dst.Name, Deep: true})
	assert.NoError(t, err)
	assert.Equal(t, 1, repo.NumPolygonsAdded)
}
