package image_store

import (
	"bytes"
	"io"
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

type TestingTransactor struct {
	Repos
}

func (m *TestingTransactor) RunInTx(
	fn func(Repos) error,
) error {
	return fn(m.Repos)
}

func Setup() (ImageStore, clc.Collection, im.Image, []byte) {
	collection := clc.NewCollection(clc.NewCollectionId(), "the-collection")
	image := im.NewImage(im.NewImageId(), collection)
	meta := []m.MetaData{{Key: "first-key", Value: 123}}
	specs := im.Specs{MIMEType: "image/jpeg"}
	image.Specs = specs
	image.Meta = meta
	data := []byte("test-data")
	repos := Repos{
		&fk.ImageRepo{
			ImageIsInCollection: true,
			ReturnSpecs:         &specs,
		},
		&fk.CollectionRepo{Return: collection},
		&fk.AnnotationRepo{},
		&fk.MetaDataRepo{ReturnList: meta},
	}
	itr := New(repos, &TestingTransactor{repos},
		&fk.FileStore{Data: data})
	return itr, collection, image, data
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	s, _, _, _ := Setup()
	s.CollectionRepo = &fk.CollectionRepo{ErrOnFind: e.ErrNotFound}
	_, err := s.Find(im.BaseImage{
		ImageId:    im.NewImageId(),
		Collection: "a-collection",
	})
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestErrOnFindLabelShouldFail(t *testing.T) {
	s, _, _, _ := Setup()
	s.AnnotationRepo = &fk.AnnotationRepo{ErrOnFindImageLabels: e.ErrInternal}
	_, err := s.Find(im.BaseImage{
		ImageId:    im.NewImageId(),
		Collection: "a-collection",
	})
	assert.NotNil(t, err)
}

func TestErrOnFindBoundingBoxesShouldFail(t *testing.T) {
	s, _, _, _ := Setup()
	s.AnnotationRepo = &fk.AnnotationRepo{ErrOnFindBoundingBoxes: e.ErrInternal}
	_, err := s.Find(im.BaseImage{
		ImageId:    im.NewImageId(),
		Collection: "a-collection",
	})
	assert.NotNil(t, err)
}

func TestErrOnFindPolygonsShouldFail(t *testing.T) {
	s, _, _, _ := Setup()
	s.AnnotationRepo = &fk.AnnotationRepo{ErrOnFindPolygons: e.ErrInternal}
	_, err := s.Find(im.BaseImage{
		ImageId:    im.NewImageId(),
		Collection: "a-collection",
	})
	assert.NotNil(t, err)
}

func TestErrOnExistsShouldFail(t *testing.T) {
	s, _, _, _ := Setup()
	s.ImageRepo = &fk.ImageRepo{ErrOnImageExistsInCollection: e.ErrInternal}
	_, err := s.Find(im.BaseImage{
		ImageId:    im.NewImageId(),
		Collection: "a-collection",
	})
	assert.NotNil(t, err)
}

func TestFindImageGivesCorrectAnnotations(t *testing.T) {
	s, collection, _, _ := Setup()
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	labels := []a.ImageLabel{{Id: a.NewAnnotationId(), Label: label}}
	bboxes := []a.BoundingBox{{Id: a.NewAnnotationId(), Label: label}}
	polygons := []a.Polygon{{Id: a.NewAnnotationId(), Label: label}}
	s.AnnotationRepo = &fk.AnnotationRepo{
		Labels:        labels,
		BoundingBoxes: bboxes, Polygons: polygons,
	}
	s.FileStore = &fk.FileStore{Data: []byte("test-data")}

	image, err := s.Find(im.BaseImage{
		ImageId:    im.NewImageId(),
		Collection: collection.Name,
	})
	assert.NoError(t, err)
	assert.Equal(t, collection.Name, image.Collection.Name)
	assert.Equal(t, 1, len(image.Labels))
	assert.Equal(t, 1, len(image.BoundingBoxes))
	assert.Equal(t, 1, len(image.Polygons))
}

func TestImageReaderGivesCorrectBytes(t *testing.T) {
	s, _, _, data := Setup()
	image, _ := s.Find(im.BaseImage{
		ImageId:    im.NewImageId(),
		Collection: "the-collection",
	})
	gotBytes, _ := io.ReadAll(image.Reader)
	assert.Equal(t, true, bytes.Equal(gotBytes, data))
}

func TestRetrieveSpecs(t *testing.T) {
	s, _, image, _ := Setup()
	r, _ := s.Find(im.BaseImage{
		ImageId:    im.NewImageId(),
		Collection: "the-collection",
	})
	assert.Equal(t, image.Specs, r.Specs)
}

func TestRetrieveMetaData(t *testing.T) {
	s, collection, image, _ := Setup()
	r, _ := s.Find(im.BaseImage{
		ImageId:    image.Id,
		Collection: collection.Name,
	})
	assert.Equal(t, image.Meta, r.Meta)
}

func SetupCopy() (ImageStore, clc.Collection, im.Image, clc.Collection, *fk.ImageRepo, *fk.AnnotationRepo) {
	srcCollection := clc.NewCollection(clc.NewCollectionId(), "src")
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst")
	image := im.NewImage(im.NewImageId(), srcCollection)
	image.AddLabel(lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	imrepo := fk.ImageRepo{
		IterateBaseImages:   []im.BaseImage{{ImageId: image.Id, Collection: srcCollection.Name}},
		ImageIsInCollection: true,
		ReturnSpecs:         &im.Specs{MIMEType: "image/jpeg"},
	}
	anrepo := fk.AnnotationRepo{Labels: image.Labels}
	repos := Repos{
		ImageRepo:      &imrepo,
		CollectionRepo: &fk.CollectionRepo{},
		AnnotationRepo: &anrepo,
		MetaRepo:       &fk.MetaDataRepo{},
	}
	transactor := TestingTransactor{repos}
	s := New(repos, &transactor, &fk.FileStore{})
	return s, srcCollection, image, dstCollection, &imrepo, &anrepo
}

func TestShallowCopyOneImage(t *testing.T) {
	store, srcCollection, image, dstCollection, imrepo, _ := SetupCopy()
	err := store.Copy(srcCollection.Name, image.Id, dstCollection.Name, false)
	assert.NoError(t, err)
	assert.Equal(t, dstCollection.Name, *imrepo.AddedIntoCollection)
}

func TestDeepCopyOneImage(t *testing.T) {
	store, srcCollection, image, dstCollection, _, anrepo := SetupCopy()
	err := store.Copy(srcCollection.Name, image.Id, dstCollection.Name, true)
	assert.NoError(t, err)
	assert.NotNil(t, anrepo.AddedAnnotationId)
}
