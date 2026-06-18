package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	presenter "github.com/lejeunel/go-image-annotator/adapters/api/json/image"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	rd "github.com/lejeunel/go-image-annotator/modules/reader"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/read"
)

func (s *Server) IngestImage(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewImage](w, r)
	if !ok {
		return
	}

	s.Image.Ingest.Execute(r.Context(), NewIngestImageRequest(*body, s.Image.AllowedImageFormats),
		presenter.NewIngestPresenter(w, s.Logger))
}

func (s *Server) ReadImage(w http.ResponseWriter, r *http.Request, collectionName, imageId string) {
	s.Image.Read.Execute(read.Request{ImageId: imageId, Collection: collectionName},
		presenter.NewReadMetaPresenter(w, s.Logger))
}

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request, params ListImagesParams) {
	req := list.Request{
		FilteringParams: im.FilteringParams{
			Collection: params.Collection,
			PageSize:   s.Image.DefaultPageSize,
		},
		OrderingParams: im.OrderingParams{IngestTime: true}}
	if p := params.Page; p != nil {
		req.Page = *p
	}
	if p := params.PageSize; p != nil {
		req.PageSize = *p
	}
	s.Image.List.Execute(req, presenter.NewListPresenter(w, s.Logger))
}

func NewIngestImageRequest(req models.NewImage, allowedImageFormats []string) ingest.Request {

	ingestReq := ingest.Request{Collection: req.Collection, Reader: rd.NewBase64ImageDecoder(allowedImageFormats, req.Data)}
	appendLabelsToIngestImageRequest(&ingestReq, req.Labels)
	appendBoundingBoxesToIngestImageRequest(&ingestReq, req.BoundingBoxes)
	return ingestReq
}

func appendBoundingBoxesToIngestImageRequest(req *ingest.Request, boxes *[]models.NewBoundingBox) {
	if boxes != nil {
		for _, box := range *boxes {
			req.BoundingBoxes = append(req.BoundingBoxes,
				ingest.BoundingBoxRequest{Xc: box.Xc, Yc: box.Yc,
					Width: box.Width, Height: box.Height})
		}
	}
}

func appendLabelsToIngestImageRequest(req *ingest.Request, labels *[]string) {
	if labels != nil {
		req.Labels = *labels
	}
}
