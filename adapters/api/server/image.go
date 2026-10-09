package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	apijson "github.com/lejeunel/go-image-annotator/adapters/api/json"
	presenter "github.com/lejeunel/go-image-annotator/adapters/api/json/image"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	ig "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/slice"
)

func (s *Server) IngestImage(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		apijson.WriteError(w, http.StatusBadRequest, "invalid multipart body")
		return
	}

	var meta models.NewImage
	var imageReader io.Reader

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			apijson.WriteError(w, http.StatusBadRequest, "error reading multipart body")
			return
		}

		switch part.FormName() {
		case "metadata":
			if err := json.NewDecoder(part).Decode(&meta); err != nil {
				apijson.WriteError(w, http.StatusBadRequest, "invalid metadata payload")
				return
			}
		case "image":
			buf, err := io.ReadAll(part) // or stream directly to your storage/hasher
			if err != nil {
				apijson.WriteError(w, http.StatusBadRequest, "error reading image data")
				return
			}
			imageReader = bytes.NewReader(buf)
		}
	}

	if imageReader == nil {
		apijson.WriteError(w, http.StatusBadRequest, "missing image part")
		return
	}
	s.Image.Ingest.Execute(r.Context(), NewImageIngestRequest(meta, imageReader),
		presenter.NewIngestPresenter(w, s.Logger))
}

func (s *Server) ReadRawImage(w http.ResponseWriter, r *http.Request, imageId string) {
	s.Image.Raw.Execute(imageId, presenter.NewRawImagePresenter(w, s.Logger))
}

func (s *Server) ReadImage(w http.ResponseWriter, r *http.Request, collectionName, imageId string) {
	s.Image.Find.Execute(find.Request{ImageId: imageId, Collection: collectionName},
		presenter.NewReadMetaPresenter(w, s.Logger))
}

func (s *Server) DeleteImage(
	w http.ResponseWriter,
	r *http.Request,
	collectionName, imageId string,
) {
	s.Image.Delete.Execute(
		r.Context(),
		delete.Request{ImageId: imageId, Collection: collectionName},
		presenter.NewDeleteImagePresenter(w, s.Logger),
	)
}

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request, params ListImagesParams) {
	pagination := pa.NewPaginationParamsFromOptional(
		params.PageSize,
		params.Page,
		s.Image.DefaultPageSize,
	)
	req := slice.Request{
		PaginationParams: pagination,
	}
	if params.Filter != nil {
		req.FilterStr = *params.Filter
	}
	if params.Order != nil {
		req.OrderStr = *params.Order
	}
	s.Image.Slice.Execute(req, presenter.NewListPresenter(w, s.Logger))
}

func NewImageIngestRequest(meta models.NewImage, reader io.Reader) ig.Request {
	ingestReq := ig.Request{
		Collection: meta.Collection,
		Reader:     reader,
	}
	appendLabelsToIngestImageRequest(&ingestReq, meta.Labels)
	appendBoundingBoxesToIngestImageRequest(&ingestReq, meta.BoundingBoxes)
	return ingestReq
}

func appendBoundingBoxesToIngestImageRequest(req *ig.Request, boxes *[]models.IngestBoundingBox) {
	if boxes != nil {
		for _, box := range *boxes {
			req.BoundingBoxes = append(req.BoundingBoxes,
				an.BoundingBoxRequest{
					Xc: box.Xc, Yc: box.Yc,
					Width: box.Width, Height: box.Height,
				})
		}
	}
}

func appendLabelsToIngestImageRequest(req *ig.Request, labels *[]string) {
	if labels != nil {
		req.Labels = *labels
	}
}
