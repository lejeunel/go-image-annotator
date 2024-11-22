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

func (s *ImageService) Delete(ctx context.Context, image *m.Image) error {
	if err := g.CheckAuthorization(ctx, "admin"); err != nil {
		return err
	}
	return s.ImageRepo.Delete(ctx, image)
}

func (s *ImageService) Save(ctx context.Context, image *m.Image) (*m.Image, error) {

	if err := g.CheckAuthorization(ctx, "im-contrib"); err != nil {
		return nil, err
	}

	if err := s.checkSHA256(image); err != nil {
		return nil, err
	}
	im, format, err := i.Decode(bytes.NewBuffer(image.Data))

	if err != nil {
		return nil, err
	}

	image.Id = uuid.New()
	image.Width = im.Bounds().Dx()
	image.Height = im.Bounds().Dy()
	image.MIMEType = "image/" + format
	s.setURI(image)

	image, err = s.ImageRepo.Create(ctx, image)
	if err != nil {
		return nil, err
	}

	if err = s.KeyValueStoreClient.Upload(ctx, image.Uri, image.Data, image.SHA256); err != nil {
		return nil, err
	}

	return image, nil

}

func (s *ImageService) GetOne(ctx context.Context, id string, withData bool) (*m.Image, error) {

	image, err := s.ImageRepo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	annotations, err := s.LabelRepo.GetAnnotationsOfImage(ctx, image)
	if err != nil {
		return nil, err
	}
	image.Annotations = annotations

	polygons, err := s.LabelRepo.GetPolygonsOfImage(ctx, image)
	if err != nil {
		return nil, err
	}
	image.Polygons = polygons

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
	pagination g.PaginationParams,
	filters *g.ImageFilterArgs,
	withData bool) ([]m.Image, *g.PaginationMeta, error) {

	if pagination.PageSize > s.MaxPageSize {
		pagination.PageSize = s.MaxPageSize
	}

	p := s.ImageRepo.Paginate(pagination.PageSize, filters)
	p.SetPage(int(pagination.Page))

	var images []m.Image
	err := p.Results(&images)

	if err != nil {
		return nil, nil, err
	}

	paginationMeta := g.NewPaginationMeta(p)
	return images, &paginationMeta, nil

}
