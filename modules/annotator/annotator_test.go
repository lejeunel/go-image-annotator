package annotator

import (
	"testing"

	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	del "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
	"github.com/stretchr/testify/assert"
)

func createAnnotator() (*Annotator, *im.Image, *FakeScroller) {
	scroller := &FakeScroller{}
	image := im.NewImage(im.NewImageId(),
		clc.NewCollection(clc.NewCollectionId(), "name"))
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	box := an.NewBoundingBox(an.NewAnnotationId(), 1, 1, 1, 1, label)
	image.AddBoundingBox(box)
	image.AddLabel(label)
	annotator := NewAnnotator(
		scroller,
		&FakeImageReader{},
		&FakeBoxAdder{},
		&FakeBoxUpdater{},
		&FakePolygonAdder{},
		&FakePolygonUpdater{},
		&FakeAnnotationDeleter{},
		&FakeLabelFetcher{},
		&FakeLabelUpdater{},
		&FakeLabelAdder{})
	return &annotator, &image, scroller

}
func TestInitializeScrollerOnStart(t *testing.T) {
	a, image, scroller := createAnnotator()
	p := &FakeScrollerPresenter{}
	a.Init(t.Context(), image.Id.String(),
		"a-collection", &FakeImageReadPresenter{}, &FakeLabelFetchPresenter{}, p)
	assert.True(t, scroller.IsInit)
}
func TestFetchLabelsOnInit(t *testing.T) {
	a, image, _ := createAnnotator()
	lp := FakeLabelFetchPresenter{}
	a.Init(t.Context(), image.Id.String(),
		"a-collection", &FakeImageReadPresenter{}, &lp, &FakeScrollerPresenter{})
	assert.NotNil(t, lp.Called)
}

func TestDrawImageOnInit(t *testing.T) {
	a, image, _ := createAnnotator()
	ip := &FakeImageReadPresenter{}
	a.Init(t.Context(), image.Id.String(),
		"a-collection", ip, &FakeLabelFetchPresenter{}, &FakeScrollerPresenter{})
	assert.True(t, ip.Called)
}
func TestAddBox(t *testing.T) {
	a, _, _ := createAnnotator()
	p := &FakeAddBoxPresenter{}
	a.AddBox(t.Context(), addbox.Request{}, p)
	assert.True(t, p.Called)
}
func TestUpdateLabel(t *testing.T) {
	a, _, _ := createAnnotator()
	p := &FakeUpdateLabelPresenter{}
	a.UpdateLabel(t.Context(), updlbl.Request{}, p)
	assert.True(t, p.Called)
}
func TestDeleteAnnotation(t *testing.T) {
	a, _, _ := createAnnotator()
	p := &FakeRemoveLabelPresenter{}
	a.DeleteAnnotation(t.Context(), del.Request{}, p)
	assert.True(t, p.Called)
}
