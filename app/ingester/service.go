package ingestion

import (
	"context"
	au "datahub/app/authorizer"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	base64 "encoding/base64"
	"fmt"
	"log/slog"
	"time"
)

type Service struct {
	Images      *im.Service
	Labels      *lbl.Service
	Collections *clc.Service
	Locations   *loc.Service
	Logger      *slog.Logger
	Authorizer  *au.Authorizer
}

type BoundingBoxIngestion struct {
	Xc     float64 `json:"xc"`
	Yc     float64 `json:"yc"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Angle  float64 `json:"angle" required:"false"`
	Label  string  `json:"label"`
	Author string  `json:"author" required:"false"`
	Date   string  `json:"date" required:"false"`
}

type ImageIngestionPayload struct {
	Data          string                 `doc:"base64-encoded data" json:"data"`
	Site          *string                `doc:"name of site" json:"site_name,omitempty"`
	Camera        *string                `doc:"name of camera" json:"camera_name,omitempty"`
	Type          *string                `json:"type,omitempty"`
	MIMEType      string                 `doc:"MIMEType" json:"mimetype"`
	CapturedAt    string                 `doc:"Acquisition datetime" json:"captured_at" required:"false"`
	Group         string                 `doc:"Group" json:"group" required:"true"`
	BoundingBoxes []BoundingBoxIngestion `json:"bounding_boxes" required:"false"`
	DryRun        bool                   `json:"dry_run" required:"false" default:"false"`
}

func NewImageToIngest(base64Data, capturedAt, type_ string) (*im.Image, error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("ingesting image: %w", err)
	}
	image, err := im.New(data)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02T15:04:05.000Z"
	if capturedAt != "" {
		capturedAt, err := time.Parse(layout, capturedAt)
		if err != nil {
			return nil, fmt.Errorf("ingesting captured_at: %w", err)
		}

		image.CapturedAt = capturedAt
	}
	image.Type = type_

	return image, nil

}

func NewIngestionService(i *im.Service, c *clc.Service, lbl *lbl.Service,
	l *loc.Service, logger *slog.Logger, auth *au.Authorizer) *Service {

	return &Service{Images: i, Collections: c, Locations: l,
		Labels:     lbl,
		Logger:     logger,
		Authorizer: auth}
}
func (s *Service) ingestBoundingBoxes(ctx context.Context, image *im.Image, collection *clc.Collection, boxes []BoundingBoxIngestion) error {
	for _, b := range boxes {
		label, err := s.Labels.FindByName(ctx, b.Label)
		if err != nil {
			return fmt.Errorf("ingesting bounding box: %w", err)
		}

		bbox, err := im.NewBoundingBox(b.Xc, b.Yc, b.Height, b.Width)
		if err != nil {
			return fmt.Errorf("ingestion bouding box: %w", err)
		}
		bbox.Annotate(label)
		if err := s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image); err != nil {
			return fmt.Errorf("ingesting bounding box %+v: %w", b, err)
		}
	}
	return nil

}

func (s *Service) Ingest(ctx context.Context, collectionName string, p ImageIngestionPayload) (*im.Image, error) {

	baseErrMsg := "ingesting image"

	type_ := ""
	if p.Type != nil {
		type_ = *p.Type
	}
	image, err := NewImageToIngest(p.Data, p.CapturedAt, type_)
	if err != nil {
		return nil, err
	}
	if (p.Site != nil) && (p.Camera != nil) {
		if err := s.Images.AssignLocation(ctx, image, *p.Site, *p.Camera); err != nil {
			return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
		}

	}

	collection, err := s.Collections.FindByName(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	if err := s.Images.ChecksumAlreadyExists(image.SHA256); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	if p.DryRun == true {
		return nil, nil
	}

	if err = s.Images.Save(ctx, image, collection); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	if err = s.ingestBoundingBoxes(ctx, image, collection, p.BoundingBoxes); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	return image, nil

}
