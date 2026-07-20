package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	presenter "github.com/lejeunel/go-image-annotator/adapters/api/json/image"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	ig "github.com/lejeunel/go-image-annotator/modules/ingester"
	rd "github.com/lejeunel/go-image-annotator/modules/reader"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
)

func (s *Server) IngestImage(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewImage](w, r)
	if !ok {
		return
	}

	s.Image.Ingest.Execute(r.Context(), NewIngestImageRequest(*body, s.Image.AllowedImageFormats),
		presenter.NewIngestPresenter(w, s.Logger))
}
func (s *Server) ReadRawImage(w http.ResponseWriter, r *http.Request, imageId string) {
	s.Image.Raw.Execute(imageId, presenter.NewRawImagePresenter(w, s.Logger))
}

func (s *Server) ReadImage(w http.ResponseWriter, r *http.Request, collectionName, imageId string) {
	s.Image.Find.Execute(find.Request{ImageId: imageId, Collection: collectionName},
		presenter.NewReadMetaPresenter(w, s.Logger))
}

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request, params ListImagesParams) {
	req := list.Request{
		Filtering: im.Filtering{
			Collection: params.Collection},
		PaginationParams: pa.PaginationParams{
			PageSize: s.Image.DefaultPageSize,
		},
		Ordering: im.Ordering{IngestTime: true}}
	if p := params.Page; p != nil {
		req.Page = *p
	}
	if p := params.PageSize; p != nil {
		req.PageSize = *p
	}
	s.Image.List.Execute(req, presenter.NewListPresenter(w, s.Logger))
}

func NewIngestImageRequest(req models.NewImage, allowedImageFormats []string) ig.Request {

	ingestReq := ig.Request{Collection: req.Collection, Reader: rd.NewBase64ImageDecoder(allowedImageFormats, req.Data)}
	appendLabelsToIngestImageRequest(&ingestReq, req.Labels)
	appendBoundingBoxesToIngestImageRequest(&ingestReq, req.BoundingBoxes)
	return ingestReq
}

func appendBoundingBoxesToIngestImageRequest(req *ig.Request, boxes *[]models.NewBoundingBox) {
	if boxes != nil {
		for _, box := range *boxes {
			req.BoundingBoxes = append(req.BoundingBoxes,
				an.BoundingBoxRequest{Xc: box.Xc, Yc: box.Yc,
					Width: box.Width, Height: box.Height})
		}
	}
}

func appendLabelsToIngestImageRequest(req *ig.Request, labels *[]string) {
	if labels != nil {
		req.Labels = *labels
	}
}
