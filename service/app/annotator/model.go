package annotator

import (
	"bufio"
	"bytes"
	"context"
	a "datahub/app/authorizer"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	e "datahub/errors"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

//go:embed templates/*
var templatesFiles embed.FS

type Annotator struct {
	Labels      *lbl.Service
	Images      *im.Service
	Collections *clc.Service
	Locations   *loc.Service
	Authorizer  *a.Authorizer
	Logger      *slog.Logger
	Rescaler    *Rescaler
	Colorizer   *Colorizer
}

type AnnotatorState struct {
	ImageId          im.ImageId
	CollectionId     clc.CollectionId
	ImageMIMEType    string
	ImageType        string
	ImageWidth       int
	ImageHeight      int
	ImageCapturedAt  time.Time
	ImageCreatedAt   time.Time
	CollectionName   string
	Group            string
	SiteName         string
	Camera           string
	ImageData        []byte
	BoundingBoxes    []*BoundingBox
	AvailableLabels  []string
	CanAnnotate      bool
	NextImage        *im.Image
	PrevImage        *im.Image
	Request          AnnotatorRequest
	AvailableSites   []loc.Site
	AvailableCameras []*loc.Camera
	CurrentSite      *loc.Site
	CurrentCamera    *loc.Camera
}

func NewAnnotator(l *lbl.Service, i *im.Service,
	c *clc.Service, loc *loc.Service, auth *a.Authorizer, logger *slog.Logger,
	targetImageWidth int) *Annotator {
	annotator := &Annotator{Labels: l, Images: i, Collections: c, Locations: loc,
		Authorizer: auth, Logger: logger,
		Rescaler:  &Rescaler{TargetWidth: targetImageWidth},
		Colorizer: &Colorizer{Colors: Palette},
	}
	return annotator
}

func (a *Annotator) DeleteAnnotation(ctx context.Context, id string) error {
	if err := a.Images.Annotations.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (a *Annotator) UpsertBoundingBox(ctx context.Context, imageId im.ImageId, collectionId clc.CollectionId, b *BoundingBox) error {
	errCtx := "upserting bounding box"
	label, err := a.Labels.FindByName(ctx, b.Label)
	if err != nil {
		return fmt.Errorf("%v: fetching label %v: %w", errCtx, b.Label, err)
	}

	image, err := a.Images.Find(ctx, imageId, collectionId, im.FetchWithRawData)
	if err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}

	a.Rescaler.BackwardTransformBoundingBoxCoords(float64(image.Width), b)
	bbox, err := im.NewBoundingBox(b.Xc, b.Yc, b.Height, b.Width)
	if err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}
	bbox.Annotate(label)

	if b.Id != "" {
		id, err := uuid.Parse(b.Id)
		if err != nil {
			return fmt.Errorf("%v: got non-null id %v: %w", errCtx, b.Id, err)
		}
		bbox.Annotation.Id = id
	}

	if err := a.Images.Annotations.UpsertBoundingBox(ctx, bbox, image); err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}
	b.Id = bbox.Annotation.Id.String()

	return nil
}

func (a *Annotator) colorizeBoundingBoxes(boxes []*BoundingBox) {
	a.Colorizer.Colorize(boxes)

}
func (a *Annotator) makeBoundingBoxPayload(ctx context.Context, image *im.Image) []*BoundingBox {
	boxes := []*BoundingBox{}
	for _, b := range image.BoundingBoxes {
		bboxPayload := &BoundingBox{Id: b.Annotation.Id.String(),
			Xc: b.Coords.Xc, Yc: b.Coords.Yc, Height: b.Coords.Height,
			Width: b.Coords.Width, Angle: b.Coords.Angle,
			Label:  b.Annotation.Label.Name,
			Author: b.Annotation.AuthorEmail,
			Date:   b.Annotation.UpdatedAt.Format("2006-01-02 / 15:04"),
		}
		a.Rescaler.ForwardTransformBoundingBoxCoords(float64(image.Width), bboxPayload)
		boxes = append(boxes, bboxPayload)
	}

	a.colorizeBoundingBoxes(boxes)
	return boxes

}

func (a *Annotator) MakeState(ctx context.Context, input AnnotatorRequest) (*AnnotatorState, error) {

	image, err := a.Images.Find(ctx, input.ImageId, input.CollectionId, im.FetchWithRawData)
	if err != nil {
		return nil, err
	}

	boxes := a.makeBoundingBoxPayload(ctx, image)

	collection, err := a.Collections.Find(ctx, input.CollectionId)
	if err != nil {
		return nil, err
	}
	availableLabels, err := a.Collections.GetAvailableLabelNames(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("making annotator state: getting all labels: %w", err)
	}

	canAnnotate := true
	if err := a.Authorizer.WantToContributeAnnotations(ctx, image.Group); err != nil {
		if errors.Is(err, e.ErrEntitlement) {
			canAnnotate = false
		} else {
			return nil, fmt.Errorf("making annotator state: %w", err)
		}
	}

	scroller := NewScroller(a.Images, input.Filter, input.Ordering)

	nextImage, err := scroller.GetNextImage(ctx, image)
	if err != nil && !errors.Is(err, e.ErrNotFound) {
		return nil, fmt.Errorf("making annotator state: %w", err)
	}
	prevImage, err := scroller.GetPrevImage(ctx, image)
	if err != nil && !errors.Is(err, e.ErrNotFound) {
		return nil, fmt.Errorf("making annotator state: %w", err)
	}

	var transformedImageBuf bytes.Buffer
	transformedImageWriter := bufio.NewWriter(&transformedImageBuf)

	err = a.Rescaler.TransformImage(bytes.NewReader(image.Data), transformedImageWriter)
	if err != nil {
		return nil, fmt.Errorf("making annotator state: %w", err)
	}

	availableSites, err := a.Locations.GetAllSites(ctx)
	if err != nil {
		return nil, fmt.Errorf("making annotator state: %w", err)
	}

	var selectedSite *loc.Site
	var selectedCamera *loc.Camera
	if image.Camera != nil {
		siteId := loc.NewSiteIdFromUUID(image.Camera.Site.Id.UUID)
		selectedSite, _ = a.Locations.FindSite(ctx, *siteId)
		cameraId := loc.NewCameraIdFromUUID(image.Camera.Id.UUID)
		selectedCamera, _ = a.Locations.FindCamera(ctx, *cameraId)
	}

	if input.SiteId != nil {
		selectedSite, err = a.Locations.FindSite(ctx, *input.SiteId)
		if err != nil {
			return nil, fmt.Errorf("making annotator state: fetching site by id: %v: %w",
				*input.SiteId, err)
		}
		if input.CameraId != nil {
			selectedCamera, err = a.Locations.FindCamera(ctx, *input.CameraId)
			if err != nil {
				return nil, fmt.Errorf("making annotator state: fetching camera by id: %v: %w",
					*input.CameraId, err)
			}
		}
	}

	var availableCameras []*loc.Camera
	if selectedSite != nil {
		availableCameras, err = a.Locations.ListCamerasOfSite(ctx, selectedSite)
		if err != nil {
			fmt.Errorf("making annotator state: fetching camera by id: %v: %w",
				*input.CameraId, err)
		}
	}

	return &AnnotatorState{
		ImageId:          image.Id,
		ImageMIMEType:    image.MIMEType,
		ImageType:        image.Type,
		ImageWidth:       image.Width,
		ImageHeight:      image.Height,
		ImageCapturedAt:  image.CapturedAt,
		ImageCreatedAt:   image.CreatedAt,
		ImageData:        transformedImageBuf.Bytes(),
		CollectionId:     image.Collection.Id,
		CollectionName:   image.Collection.Name,
		Group:            image.Group,
		BoundingBoxes:    boxes,
		SiteName:         image.GetSiteName(),
		Camera:           image.GetCameraName(),
		CanAnnotate:      canAnnotate,
		NextImage:        nextImage,
		PrevImage:        prevImage,
		Request:          input,
		AvailableSites:   availableSites,
		AvailableCameras: availableCameras,
		AvailableLabels:  availableLabels,
		CurrentSite:      selectedSite,
		CurrentCamera:    selectedCamera,
	}, nil
}
