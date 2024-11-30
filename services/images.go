package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
	e "go-image-annotator/errors"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
	i "image"
	_ "image/png"
)

type ImageService struct {
	KeyValueStoreClient KeyValueStoreClient
	ImageRepo           r.ImageRepo
	LabelRepo           r.AnnotationRepo
	CollectionRepo      r.CollectionRepo
	MaxPageSize         int
	DefaultPageSize     int
	RemoteScheme        string
	RemoteBucketName    string
}

func (s *ImageService) checkSHA256(image *m.Image) error {
	if image.SHA256 == "" {
		return nil

	}

	h := sha256.New()
	h.Write(image.Data)
	checkSum := hex.EncodeToString(h.Sum(nil))

	if image.SHA256 != checkSum {
		return &e.ErrCheckSum{}

	}

	return nil

}

func (s *ImageService) setURI(image *m.Image) {
	image.Uri = s.RemoteScheme + "://" + s.RemoteBucketName + "/" + image.Id.String() + ".png"
}

func (s *ImageService) Delete(ctx context.Context, image *m.Image, collection *m.Collection) error {
	if err := g.CheckAuthorization(ctx, "admin"); err != nil {
		return err
	}
	for _, a := range image.Annotations {
		err := s.LabelRepo.DeleteAnnotation(ctx, a)
		if err != nil {
			return err
		}

	}

	return s.ImageRepo.Delete(ctx, image)
}

func (s *ImageService) Save(ctx context.Context, image *m.Image, collection *m.Collection) error {
	if err := s.saveImage(ctx, image); err != nil {
		return err
	}
	err := s.CollectionRepo.AssignImageToCollection(ctx, image, collection)
	if err != nil {
		return err
	}

	return nil

}

func (s *ImageService) saveImage(ctx context.Context, image *m.Image) error {

	if err := g.CheckAuthorization(ctx, "im-contrib"); err != nil {
		return err
	}

	if err := s.checkSHA256(image); err != nil {
		return err
	}
	im, format, err := i.Decode(bytes.NewBuffer(image.Data))

	if err != nil {
		return err
	}

	image.Id = uuid.New()
	image.Width = im.Bounds().Dx()
	image.Height = im.Bounds().Dy()
	image.MIMEType = "image/" + format
	s.setURI(image)

	image, err = s.ImageRepo.Create(ctx, image)
	if err != nil {
		return err
	}

	if err = s.KeyValueStoreClient.Upload(ctx, image.Uri, image.Data, image.SHA256); err != nil {
		return err
	}

	return nil

}

func (s *ImageService) Get(ctx context.Context, collection_id string, image_id string, withData bool) (*m.Image, error) {
	image, err := s.getBase(ctx, image_id, withData)
	if err != nil {
		return nil, err
	}

	image, err = s.prependAuxiliaryFields(ctx, image, collection_id)
	if err != nil {
		return nil, err
	}

	return image, nil

}

func (s *ImageService) prependAuxiliaryFields(ctx context.Context, image *m.Image, collection_id string) (*m.Image, error) {
	collection, err := s.CollectionRepo.Get(ctx, collection_id)
	if err != nil {
		return nil, err
	}

	annotations, err := s.LabelRepo.GetAnnotationsOfImage(ctx, image, collection)

	if err != nil {
		return nil, err
	}
	image.Annotations = annotations

	bboxes, err := s.LabelRepo.GetBoundingBoxesOfImage(ctx, image)
	if err != nil {
		return nil, err
	}
	image.BoundingBoxes = bboxes

	return image, nil

}

func (s *ImageService) getBase(ctx context.Context, id string, withData bool) (*m.Image, error) {

	image, err := s.ImageRepo.GetOne(ctx, id)

	if err != nil {
		return nil, err
	}

	if withData {
		data, err := s.KeyValueStoreClient.Download(ctx, image.Uri)

		if err != nil {
			return nil, err
		}
		image.Data = data
	}

	return image, nil
}

func (s *ImageService) GetPage(
	ctx context.Context,
	collection_id string,
	pagination g.PaginationParams,
	withData bool) ([]m.Image, *g.PaginationMeta, error) {

	p := s.ImageRepo.Paginate(pagination.PageSize, &g.ImageFilterArgs{CollectionId: collection_id})
	p.SetPage(int(pagination.Page))

	var images []m.Image
	var augmentedImages []m.Image
	err := p.Results(&images)

	// augment objects with annotations
	for _, image := range images {
		image, err := s.Get(ctx, collection_id, image.Id.String(), withData)
		if err != nil {
			return nil, nil, err
		}
		augmentedImages = append(augmentedImages, *image)

	}

	if err != nil {
		return nil, nil, err
	}

	paginationMeta := g.NewPaginationMeta(p)
	return augmentedImages, &paginationMeta, nil

}
