package images

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	url "net/url"
	"slices"
	"strings"
	"time"

	au "datahub/app/authorizer"
	clc "datahub/domain/collections"
	loc "datahub/domain/locations"
	e "datahub/errors"
	g "datahub/generic"

	clk "github.com/jonboulle/clockwork"

	"path/filepath"
)

type RawImage struct {
	Data     []byte
	MIMEType string
}

func NewImageService(kvs *KeyValueStoreClient,
	imageRepo ImageRepo,
	annotationService *AnnotationService,
	locationService *loc.Service,
	collectionService *clc.Service,
	maxPageSize int,
	logger *slog.Logger,
	allowedTypes []string,
	authorizer *au.Authorizer,
	clock clk.Clock,
) *Service {
	return &Service{
		KeyValueStoreClient: *kvs,
		Repo:                imageRepo,
		Annotations:         annotationService,
		Locations:           locationService,
		CollectionService:   collectionService,
		MaxPageSize:         maxPageSize,
		Logger:              logger,
		AllowedTypes:        allowedTypes,
		Authorizer:          authorizer,
		Clock:               clock}
}

type Service struct {
	KeyValueStoreClient KeyValueStoreClient
	Repo                ImageRepo
	CollectionService   *clc.Service
	Annotations         *AnnotationService
	Locations           *loc.Service
	Logger              *slog.Logger
	MaxPageSize         int
	AllowedTypes        []string
	Authorizer          *au.Authorizer
	Clock               clk.Clock
}

func makeURI(id ImageId, scheme, root, format string) url.URL {
	path := root + "/"

	path += id.String()
	path += "." + format
	return url.URL{Scheme: scheme, Path: path}

}

func (s *Service) DeleteCollection(ctx context.Context, collection *clc.Collection) error {
	if err := s.Authorizer.WantToDeleteCollectionOrItsContent(ctx, collection.Group); err != nil {
		return fmt.Errorf("deleting collection: %w", err)
	}

	if err := s.Annotations.AnnotationRepo.DeleteAllAnnotations(collection); err != nil {
		return fmt.Errorf("deleting annotations of collection: %w", err)
	}

	if err := s.Repo.DeleteImagesInCollection(collection); err != nil {
		return fmt.Errorf("deleting all images in collection %v: %w", collection.Id, err)
	}
	if err := s.CollectionService.Delete(ctx, collection.Id); err != nil {
		return fmt.Errorf("deleting collection %v: %w", collection.Id, err)
	}
	return nil
}

func (s *Service) Patch(ctx context.Context, imageId ImageId, patches g.JSONPatches) (*BaseImage, error) {
	baseErrMsg := "patching image"
	image, err := s.GetBase(ctx, imageId, FetchImageOptions{IncludeRawData: false})
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	if err := s.appendCamera(ctx, image); err != nil {
		return nil, fmt.Errorf("updating image: %w", err)
	}
	original := ImageUpdatables{
		CapturedAt: image.CapturedAt.Format("2006-01-02T15:04:05.000Z"),
		Site:       image.GetSiteName(),
		Camera:     image.GetCameraName(),
		Type_:      image.Type,
	}
	originalJSONBytes, err := json.Marshal(&original)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	modifiedBytes, err := patches.Apply(originalJSONBytes)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	var modified ImageUpdatables
	if err := json.Unmarshal(modifiedBytes, &modified); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	return s.Update(ctx, imageId, modified)
}

func (s *Service) Update(ctx context.Context, imageId ImageId,
	payload ImageUpdatables) (*BaseImage, error) {

	baseErrMsg := "updating image"

	base, err := s.GetBase(ctx, imageId, FetchImageOptions{IncludeRawData: false})
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	if err := s.appendCamera(ctx, base); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	err = s.Authorizer.WantToContributeImages(ctx, base.Group)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	layout := "2006-01-02T15:04:05.000Z"
	parsedCapturedAt, err := time.Parse(layout, payload.CapturedAt)
	if err != nil {
		return nil, fmt.Errorf("%v: parsing captured_at %v: %w", baseErrMsg, payload.CapturedAt, err)
	}
	if err := s.updateTypeAndTime(base, payload.Type_, parsedCapturedAt); err != nil {
		return nil, err
	}

	if err := s.AssignLocation(ctx, base.Id, payload.Site, payload.Camera); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	base, err = s.GetBase(ctx, base.Id, FetchMetaOnly)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	return base, nil
}

func (s *Service) AssignToCollection(ctx context.Context, image *Image, collection *clc.Collection) error {
	err := s.Repo.AssignToCollection(image, collection)
	if err != nil {
		return err
	}
	image.Collection.Id = collection.Id
	image.Group = collection.Group

	s.CollectionService.Touch(ctx, collection)
	return nil
}

func (s *Service) Save(ctx context.Context, image *Image, collection *clc.Collection) error {
	err := s.Authorizer.WantToContributeImages(ctx, collection.Group)
	if err != nil {
		return fmt.Errorf("adding image: %w", err)
	}

	if image.Type != "" {
		if !slices.Contains(s.AllowedTypes, image.Type) {
			return fmt.Errorf("image type %v not allowed. Allowed values are %v", image.Type, s.AllowedTypes)
		}
	}
	format := strings.Split(image.MIMEType, "/")[1]
	uri := makeURI(image.Id,
		s.KeyValueStoreClient.Scheme(),
		s.KeyValueStoreClient.Root(), format)

	image.Uri = uri
	image.FileName = filepath.Base(uri.Path)
	image.Group = collection.Group
	image.Collection = collection

	image.CreatedAt = s.Clock.Now()
	image.UpdatedAt = s.Clock.Now()

	if err := s.saveImage(ctx, image); err != nil {
		err = fmt.Errorf("saving image with URI to local store: %v: %w", image.Uri.String(), err)
		s.Logger.Error(err.Error())
		return err
	}
	if err := s.AssignToCollection(ctx, image, collection); err != nil {
		err = fmt.Errorf("saving image with URI: %v: assigning to collection: %v: %w",
			image.Uri.String(), collection.Name, err)
		s.Logger.Error(err.Error())
		return err
	}

	return nil

}

func (s *Service) ChecksumAlreadyExists(sha256 string) error {
	images, err := s.Repo.ListWithChecksum(sha256)
	if err != nil {
		return err
	}
	matchedIds := []string{}
	if len(images) > 0 {
		for _, image := range images {
			matchedIds = append(matchedIds, image.Id.String())
		}
		return fmt.Errorf("checking for duplicate image using SHA256: found matching id(s) %v: %w",
			strings.Join(matchedIds, ", "), e.ErrDuplication)

	}
	return nil
}

func (s *Service) Find(ctx context.Context, imageId ImageId, collectionId clc.CollectionId, opts FetchImageOptions) (*Image, error) {
	baseErrMsg := "getting image by id and collection id"
	base, err := s.GetBase(ctx, imageId, opts)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	image, err := s.appendAggregateResources(ctx, base, collectionId, opts)
	if err != nil {
		return nil, err
	}

	return image, nil

}
func (s *Service) appendCamera(ctx context.Context, image *BaseImage) error {
	if image.CameraId != nil {
		camera, err := s.Locations.FindCamera(ctx, *image.CameraId)
		if err != nil {
			return fmt.Errorf("appending camera to image: %w", err)
		}
		image.Camera = camera
	}
	return nil

}

func (s *Service) appendAggregateResources(ctx context.Context, base *BaseImage, collectionId clc.CollectionId, opts FetchImageOptions) (*Image, error) {

	uri := makeURI(base.Id,
		s.KeyValueStoreClient.Scheme(),
		s.KeyValueStoreClient.Root(),
		strings.Split(base.MIMEType, "/")[1])

	base.Uri = uri

	baseErrMsg := "appending aggregate resources"
	if err := s.appendCamera(ctx, base); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	collection, err := s.CollectionService.Find(ctx, collectionId)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	image := &Image{Id: base.Id, FileName: base.FileName, CapturedAt: base.CapturedAt, CreatedAt: base.CreatedAt, UpdatedAt: base.UpdatedAt,
		SHA256: base.SHA256, MIMEType: base.MIMEType, Width: base.Width, Height: base.Height, Type: base.Type, Collection: collection, Camera: base.Camera,
		Uri: base.Uri, Group: base.Group}

	if err := s.appendDataAndAnnotations(ctx, image, opts); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	return image, nil

}
func (s *Service) GetBase(ctx context.Context, imageId ImageId, opts FetchImageOptions) (*BaseImage, error) {

	errCtx := "fetching image by id"
	image, err := s.Repo.GetBase(imageId)
	if err != nil {
		return nil, fmt.Errorf("%v: %v: %w", errCtx, imageId, err)
	}
	if image.CameraId != nil {
		camera, err := s.Locations.FindCamera(ctx, *image.CameraId)
		if err != nil {
			return nil, fmt.Errorf("%v: %v: %w", errCtx, imageId, err)
		}
		image.Camera = camera
	}

	return image, nil
}

func (s *Service) GetRaw(ctx context.Context, imageId *ImageId) (*RawImage, error) {

	baseErrMsg := fmt.Sprintf("fetching raw image data by id %v", imageId)
	base, err := s.Repo.GetBase(*imageId)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	data, err := s.getImageData(ctx, base.FileName)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	return &RawImage{Data: data, MIMEType: base.MIMEType}, nil

}

func (s *Service) NumImages(ctx context.Context, filters FilterArgs) (int64, error) {

	return s.Repo.Count(filters)
}

func (s *Service) GetAdjacent(ctx context.Context, currentImage *Image, filters FilterArgs,
	ordering OrderingArgs, previous bool, opts FetchImageOptions) (*Image, error) {
	base, err := s.Repo.GetAdjacent(currentImage, filters, ordering, previous)
	if err != nil {
		return nil, fmt.Errorf("getting adjacent image: %w", err)
	}
	image, err := s.appendAggregateResources(ctx, base, currentImage.Collection.Id, opts)
	if err != nil {
		return nil, fmt.Errorf("getting adjacent image: %w", err)
	}
	return image, nil
}

func (s *Service) List(
	ctx context.Context,
	filters FilterArgs,
	orderings OrderingArgs,
	pagination g.PaginationParams,
	opts FetchImageOptions) ([]Image, *g.PaginationMeta, error) {

	baseErrMsg := "listing images"
	if err := pagination.Validate(s.MaxPageSize); err != nil {
		return nil, nil, fmt.Errorf("%v: validating page size %v: %w",
			s.MaxPageSize, baseErrMsg, err)
	}

	imageBaseList, paginationMeta, err := s.Repo.List(filters, orderings, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("%v: listing images: %w", baseErrMsg, err)
	}

	images := []Image{}
	for i := 0; i < len(imageBaseList); i++ {
		base := &imageBaseList[i]
		image, err := s.appendAggregateResources(ctx, base, base.CollectionId, opts)
		images = append(images, *image)
		if err != nil {
			return nil, nil, err
		}
	}

	return images, paginationMeta, nil

}

func (s *Service) ImportImage(ctx context.Context, sourceImage *Image,
	destinationCollectionId clc.CollectionId, opts ImportImageOptions) error {

	baseErrMsg := fmt.Sprintf("importing image %v from collection %v to collection %v",
		sourceImage.Id, sourceImage.Collection.Id, destinationCollectionId)

	destinationCollection, err := s.CollectionService.Find(ctx, destinationCollectionId)
	if err != nil {
		return fmt.Errorf("%v: Fetching destination collection: %w", baseErrMsg, err)
	}

	imageFoundInDestination, err := s.Repo.ImageIsInCollection(sourceImage, destinationCollection)
	if err != nil {
		return fmt.Errorf("%v: checking whether source image exist in destination collection: %w",
			baseErrMsg, err)
	}

	if imageFoundInDestination == true {
		return fmt.Errorf("%v: source image already exists in destination collection: %w",
			baseErrMsg, e.ErrDuplication)

	}

	if err := s.Repo.AssignToCollection(sourceImage, destinationCollection); err != nil {
		return fmt.Errorf("%v: assigning base image meta-data: %w", baseErrMsg, err)
	}

	destinationImage, _ := s.Find(ctx, sourceImage.Id, destinationCollection.Id, FetchMetaOnly)
	if opts.ImportAnnotations {
		for _, b := range sourceImage.BoundingBoxes {
			if err := s.Annotations.applyBoundingBox(ctx, b, destinationImage); err != nil {
				return fmt.Errorf("%v: copying bounding box when importing: %w", baseErrMsg, err)
			}
		}
	}

	return nil
}

func (s *Service) RemoveFromCollection(ctx context.Context, image *Image) error {

	baseErrMsg := fmt.Sprintf("removing image %v from collection %v", image.Id, image.Collection.Id)
	if err := s.Authorizer.WantToDeleteCollectionOrItsContent(ctx, image.Group); err != nil {
		return fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	for _, a := range image.Annotations {
		err := s.Annotations.Delete(ctx, a.Id.String())
		if err != nil {
			return fmt.Errorf("%v: %w", baseErrMsg, err)
		}

	}
	return s.Repo.RemoveImageFromCollection(image)
}

func (s *Service) appendDataAndAnnotations(ctx context.Context, image *Image, opts FetchImageOptions) error {
	if err := s.Annotations.AppendAnnotations(ctx, image); err != nil {
		return fmt.Errorf("appending annotations: %w", err)
	}
	if opts.IncludeRawData == true {
		if err := s.appendImageData(ctx, image); err != nil {
			return fmt.Errorf("appending raw data: %w", err)
		}
	}

	return nil

}

func (s *Service) getImageData(ctx context.Context, filename string) ([]byte, error) {
	uri := url.URL{Scheme: s.KeyValueStoreClient.Scheme(), Path: strings.Join([]string{s.KeyValueStoreClient.Root(),
		filename}, "/")}
	return s.KeyValueStoreClient.Download(ctx, uri.String())

}

func (s *Service) appendImageData(ctx context.Context, image *Image) error {
	data, err := s.getImageData(ctx, image.FileName)

	if err != nil {
		return fmt.Errorf("appending image data: %w", err)
	}
	image.Data = data

	return nil

}

func (s *Service) saveImage(ctx context.Context, image *Image) error {

	image, err := s.Repo.Create(image)
	if err != nil {
		return err
	}

	s.Logger.Debug(fmt.Sprintf("uploading to %v", image.Uri))
	if err = s.KeyValueStoreClient.Upload(ctx, image.Uri.String(), image.Data, image.SHA256); err != nil {
		return err
	}

	return nil

}

func (s *Service) updateTypeAndTime(image *BaseImage, type_ string, capturedAt time.Time) error {
	if type_ != "" {
		if !slices.Contains(s.AllowedTypes, type_) {
			return fmt.Errorf("updating image: image type %v not allowed. Allowed values are %v", type_, s.AllowedTypes)
		}
	}

	if err := s.Repo.Update(image.Id, type_, capturedAt); err != nil {
		return fmt.Errorf("updating image fields: %w", err)
	}
	return nil
}

func (s *Service) AssignLocation(ctx context.Context, image_id ImageId, siteName, cameraName string) error {
	baseErrMsg := fmt.Sprintf("assigning location on image %v with site '%v' and camera '%v'", image_id, siteName, cameraName)
	site, err := s.Locations.FindSiteByName(ctx, siteName)
	if err != nil {
		return fmt.Errorf("%v: fetching site by name: %w", baseErrMsg, err)
	}
	camera, err := s.Locations.FindCameraByName(ctx, site, cameraName)
	if err != nil {
		return fmt.Errorf("%v: fetching camera by name: %w", baseErrMsg, err)
	}
	if err := s.AssignCamera(ctx, camera.Id, image_id); err != nil {
		return fmt.Errorf("%v: assigning camera: %w", baseErrMsg, err)

	}
	return nil

}

func (s *Service) AssignCamera(ctx context.Context, camera_id loc.CameraId, image_id ImageId) error {
	_, err := s.Locations.FindCamera(ctx, camera_id)
	if err != nil {
		return fmt.Errorf("assigning camera to site: %w", err)
	}

	if err := s.Repo.AssignCamera(camera_id, image_id); err != nil {
		return err
	}

	return nil
}

func (s *Service) UnassignCamera(ctx context.Context, id ImageId) error {
	return s.Repo.UnassignCamera(id)

	return nil
}
