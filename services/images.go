package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
	i "image"
	_ "image/png"
)

type ImageService struct {
	KeyValueStoreClient KeyValueStoreClient
	ImageRepo           r.ImageRepo
	LabelRepo           r.LabelRepo
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
	return s.ImageRepo.Delete(ctx, image)
}

func (s *ImageService) Save(ctx context.Context, image *m.Image) (*m.Image, error) {

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

func (s *ImageService) GetOne(ctx context.Context, id string) (*m.Image, error) {

	image, err := s.ImageRepo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	labels, err := s.LabelRepo.GetLabelsOfImage(ctx, image)
	if err != nil {
		return nil, err
	}
	image.Labels = labels

	polygons, err := s.LabelRepo.GetPolygonsOfImage(ctx, image)
	if err != nil {
		return nil, err
	}
	image.Polygons = polygons

	data, err := s.KeyValueStoreClient.Download(ctx, image.Uri)

	if err != nil {
		return nil, err
	}

	image.Data = data

	return image, nil
}
