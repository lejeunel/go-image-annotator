package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	presenter "github.com/lejeunel/go-image-annotator/adapters/api/json/image"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	ig "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
)

func (s *Server) IngestImage(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "invalid multipart body", http.StatusBadRequest)
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
			http.Error(w, "error reading multipart part", http.StatusBadRequest)
			return
		}

		switch part.FormName() {
		case "metadata":
			if err := json.NewDecoder(part).Decode(&meta); err != nil {
				http.Error(w, "invalid metadata json", http.StatusBadRequest)
				return
			}
		case "image":
			buf, err := io.ReadAll(part) // or stream directly to your storage/hasher
			if err != nil {
				http.Error(w, "error reading image data", http.StatusInternalServerError)
				return
			}
			imageReader = bytes.NewReader(buf)
		}
	}

	if imageReader == nil {
		http.Error(w, "missing image part", http.StatusBadRequest)
		return
	}
	s.Image.Ingest.Execute(r.Context(), NewIngestImageRequest(meta, imageReader),
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
		PaginationParams: pa.PaginationParams{
			PageSize: *params.PageSize,
			Page:     *params.Page,
		},
	}
	if params.Filter != nil {
		req.FilterStr = *params.Filter
	}
	if params.Order != nil {
		req.OrderStr = *params.Order
	}
	s.Image.List.Execute(req, presenter.NewListPresenter(w, s.Logger))
}

func NewIngestImageRequest(meta models.NewImage, reader io.Reader) ig.Request {
	ingestReq := ig.Request{
		Collection: meta.Collection,
		Reader:     reader,
	}
	appendLabelsToIngestImageRequest(&ingestReq, meta.Labels)
	appendBoundingBoxesToIngestImageRequest(&ingestReq, meta.BoundingBoxes)
	return ingestReq
}

func appendBoundingBoxesToIngestImageRequest(req *ig.Request, boxes *[]models.NewBoundingBox) {
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
