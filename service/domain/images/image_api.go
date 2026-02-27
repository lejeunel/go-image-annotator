package images

import (
	clc "datahub/domain/collections"
	loc "datahub/domain/locations"
	e "datahub/errors"
	g "datahub/generic"

	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type AnnotationUpdateRequest struct {
	Id   uuid.UUID `query:"id"`
	Body struct {
		Label string `json:"label"`
	}
}

type AnnotationResponse struct {
	Id    uuid.UUID `json:"id"`
	Label string    `json:"label"`
}

type ImageToImport struct {
	ImageId              string `doc:"ID of image" json:"image_id"`
	SourceCollectionName string `doc:"Name of source collection" json:"collection_name"`
}

type CollectionImportRequest struct {
	DestinationCollectionName string `path:"name" doc:"destination collection name"`
	WithAnnotations           bool   `doc:"Also import annotations" query:"deep"`
	Body                      struct {
		Sources []ImageToImport `doc:"Specification of image to import" json:"source"`
	}
}

type GetRawImageRequest struct {
	Id string `path:"id"`
}

type GetBaseImageRequest struct {
	Id string `path:"id"`
}

type DeleteImageRequest struct {
	CollectionName string `query:"collection_name" doc:"Name of collection"`
	ImageId        string `query:"image_id" doc:"Id of image"`
}

type DeleteCollectionRequest struct {
	CollectionId string `path:"id" doc:"Id of collection"`
}

type ImageRequest struct {
	ImageId        string `path:"id" doc:"Id of image"`
	CollectionName string `path:"collection_name" doc:"Name of collection"`
}

type ImagesRequest struct {
	g.PaginationParams
	CollectionId   string `query:"collection_id" required:"false" doc:"Id of collection"`
	CollectionName string `query:"collection_name" required:"false" doc:"Name of collection"`
	CameraId       string `query:"camera_id" required:"false" doc:"Id of camera"`
	LabelId        string `query:"label_id" required:"false" doc:"Id of camera"`
}

type ImageUpdateRequest struct {
	Id   string `path:"id"`
	Body ImageUpdatables
}

type ImagePatchRequest struct {
	Id   string `path:"id"`
	Body g.JSONPatches
}

type ImagesResponse struct {
	Pagination g.PaginationMeta `json:"pagination"`
	Images     []ImageResponse  `json:"images"`
}

type ImageInput struct {
	Body struct {
		CollectionId string `doc:"Collection ID" json:"collection_id"`
		Data         string `doc:"Image data (base64)" json:"data"`
		MIMEType     string `doc:"MIMEType" json:"mimetype"`
		CapturedAt   string `doc:"Acquisition datetime" json:"captured_at" required:"false"`
		Type         string `doc:"Image type" json:"type" required:false`
	}
}

type Labels []string
type BoundingBoxResponse struct {
	Id     uuid.UUID `json:"id"`
	Angle  float64   `json:"angle"`
	Xc     float64   `json:"xc"`
	Yc     float64   `json:"yc"`
	Width  float64   `json:"width"`
	Height float64   `json:"height"`
	Labels Labels    `json:"labels"`
}

type ImageBaseResponseBody struct {
	Site        string `json:"site"`
	Camera      string `json:"camera"`
	CapturedAt  string `json:"captured_at"`
	Type        string `json:"type"`
	Transmitter string `json:"transmitter"`
}
type ImageBaseResponse struct {
	Body ImageBaseResponseBody
}

type ImageResponse struct {
	Id             string                `json:"id"`
	Site           string                `json:"site,omitempty"`
	Camera         string                `json:"camera,omitempty"`
	CapturedAt     string                `json:"captured_at"`
	CollectionId   string                `json:"collection_id"`
	CollectionName string                `json:"collection_name"`
	Data           []byte                `json:"data,omitempty"`
	MIMEType       string                `json:"mimetype"`
	Labels         Labels                `json:"labels"`
	Transmitter    string                `json:"transmitter"`
	BoundingBoxes  []BoundingBoxResponse `json:"bounding_boxes"`
	RawURL         string                `json:"raw_url"`
}

type ImageGetResponse struct {
	Body ImagesResponse
}

func NewImageBaseResponse(image *BaseImage) *ImageBaseResponse {
	return &ImageBaseResponse{Body: ImageBaseResponseBody{Site: image.GetSiteName(),
		Camera: image.GetCameraName(), CapturedAt: image.CapturedAt.Format("2006-01-02T15:04:05.000Z"),
		Type: image.Type, Transmitter: image.GetTransmitter()}}
}

func NewImageResponse(image *Image, rawURLbuilder *RawImageURLBuilder) ImageResponse {
	response := ImageResponse{
		Id:             image.Id.String(),
		Site:           image.GetSiteName(),
		Camera:         image.GetCameraName(),
		CollectionId:   image.Collection.Id.String(),
		CollectionName: image.Collection.Name,
		Transmitter:    image.GetTransmitter(),
		CapturedAt:     image.CapturedAt.Format("2006-01-02T15:04:05.000Z"),
		MIMEType:       image.MIMEType,
		RawURL:         rawURLbuilder.Build(image.Id),
	}

	if len(image.Annotations) > 0 {
		for _, a := range image.Annotations {
			response.Labels = append(response.Labels, a.Label.Name)
		}
	}

	if len(image.BoundingBoxes) > 0 {
		for _, bbox := range image.BoundingBoxes {
			bbox := BoundingBoxResponse{
				Id:     bbox.Annotation.Id,
				Angle:  0,
				Xc:     bbox.Coords.Xc,
				Yc:     bbox.Coords.Yc,
				Width:  bbox.Coords.Width,
				Height: bbox.Coords.Height,
				Labels: strings.Split(bbox.Annotation.Label.String(), ",")}
			response.BoundingBoxes = append(response.BoundingBoxes, bbox)
		}
	}

	return response
}

type RawImageURLBuilder struct {
	BaseAPIURLBuilder *g.APIURLBuilder
	Endpoint          string
}

func NewRawImageURLBuilder(baseBuilder *g.APIURLBuilder, endpoint string) *RawImageURLBuilder {
	return &RawImageURLBuilder{BaseAPIURLBuilder: baseBuilder,
		Endpoint: endpoint}
}

func (b *RawImageURLBuilder) Build(id ImageId) string {
	return fmt.Sprintf("%v/%v", b.BaseAPIURLBuilder.Build(b.Endpoint), id.String())
}

type ImagesHTTPController struct {
	CollectionService  *clc.Service
	LocationService    *loc.Service
	Images             *Service
	RawImageURLBuilder *RawImageURLBuilder
}

func NewImagesHTTPController(imageService *Service,
	collectionService *clc.Service,
	locationService *loc.Service,
	rawImageURLBuilder *RawImageURLBuilder) *ImagesHTTPController {
	return &ImagesHTTPController{Images: imageService,
		CollectionService:  collectionService,
		LocationService:    locationService,
		RawImageURLBuilder: rawImageURLBuilder,
	}

}

func (h *ImagesHTTPController) GetImages(ctx context.Context, input *ImagesRequest) (*ImageGetResponse, error) {

	filters, err := NewImageFilterFromString(input.CollectionId, input.CollectionName, input.CameraId, input.LabelId)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	images, p, err := h.Images.List(ctx,
		*filters,
		OrderingArgs{},
		g.PaginationParams{Page: input.Page, PageSize: input.PageSize},
		FetchMetaOnly,
	)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	var response ImagesResponse
	response.Pagination = *p

	for _, im := range images {
		record := NewImageResponse(&im, h.RawImageURLBuilder)
		response.Images = append(response.Images, record)
	}

	return &ImageGetResponse{Body: response}, nil
}

func (h *ImagesHTTPController) Delete(ctx context.Context, input *DeleteImageRequest) (*struct{}, error) {
	imageId, err := NewImageIdFromString(input.ImageId)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	collection, err := h.CollectionService.FindByName(ctx, input.CollectionName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	image, err := h.Images.Find(ctx, *imageId, collection.Id, FetchMetaOnly)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	err = h.Images.RemoveFromCollection(ctx, image)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return nil, nil

}

func (h *ImagesHTTPController) DeleteAll(ctx context.Context, input *DeleteCollectionRequest) (*struct{}, error) {

	collectionId, err := clc.NewCollectionIdFromString(input.CollectionId)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	collection, err := h.CollectionService.Find(ctx, *collectionId)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	if err := h.Images.DeleteCollection(ctx, collection); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return nil, nil

}

func (h *ImagesHTTPController) GetOne(ctx context.Context, input *ImageRequest) (*ImageResponse, error) {
	imageId, err := NewImageIdFromString(input.ImageId)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	collection, err := h.CollectionService.FindByName(ctx, input.CollectionName)

	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	image, err := h.Images.Find(ctx, *imageId, collection.Id,
		FetchMetaOnly,
	)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	response := NewImageResponse(image, h.RawImageURLBuilder)
	return &response, nil
}

func (h *ImagesHTTPController) GetBase(ctx context.Context, input *GetBaseImageRequest) (*ImageBaseResponse, error) {
	imageId, err := NewImageIdFromString(input.Id)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	image, err := h.Images.GetBase(ctx, *imageId, FetchMetaOnly)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return NewImageBaseResponse(image), nil
}

func (h *ImagesHTTPController) extractUUIDFromPath(path string) string {
	segments := strings.Split(path, "/")
	return segments[len(segments)-1]
}

func (h *ImagesHTTPController) GetRaw(w http.ResponseWriter, r *http.Request) {

	imageId, err := NewImageIdFromString(h.extractUUIDFromPath(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), e.ToHTTPStatus(err))
		return
	}
	rawImage, err := h.Images.GetRaw(r.Context(), imageId)
	if err != nil {
		http.Error(w, err.Error(), e.ToHTTPStatus(err))
		return
	}

	// Set headers manually for file response
	w.Header().Set("Content-Type", rawImage.MIMEType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(rawImage.Data)))
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(rawImage.Data)

}

func (h *ImagesHTTPController) Update(ctx context.Context, input *ImageUpdateRequest) (*ImageBaseResponse, error) {

	imageId, err := NewImageIdFromString(input.Id)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	image, err := h.Images.Update(ctx, *imageId, input.Body)

	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return NewImageBaseResponse(image), nil

}

func (h *ImagesHTTPController) Patch(ctx context.Context, input *ImagePatchRequest) (*ImageBaseResponse, error) {
	imageId, err := NewImageIdFromString(input.Id)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	image, err := h.Images.Patch(ctx, *imageId, input.Body)

	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return NewImageBaseResponse(image), nil

}

func (h *ImagesHTTPController) Import(ctx context.Context, input *CollectionImportRequest) (*struct{}, error) {
	destinationCollection, err := h.CollectionService.FindByName(ctx, input.DestinationCollectionName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	for _, source := range input.Body.Sources {
		sourceImageId, err := NewImageIdFromString(source.ImageId)
		if err != nil {
			return nil, e.ToHumaStatusError(err)
		}

		sourceCollection, err := h.CollectionService.FindByName(ctx, source.SourceCollectionName)
		if err != nil {
			return nil, e.ToHumaStatusError(err)
		}

		srcImage, err := h.Images.Find(ctx, *sourceImageId, sourceCollection.Id, FetchMetaOnly)

		importOpts := ImportImageWithoutAnnotations
		if input.WithAnnotations == true {
			importOpts = ImportImageWithAnnotations
		}
		if err := h.Images.ImportImage(ctx, srcImage, destinationCollection.Id, importOpts); err != nil {
			return nil, e.ToHumaStatusError(err)
		}
	}

	return nil, nil

}

func (h *ImagesHTTPController) UpdateAnnotation(ctx context.Context, input *AnnotationUpdateRequest) (*AnnotationResponse, error) {
	if err := h.Images.Annotations.UpdateLabel(ctx, input.Id.String(), input.Body.Label); err != nil {
		return nil, err
	}

	return &AnnotationResponse{Id: input.Id, Label: input.Body.Label}, nil

}
