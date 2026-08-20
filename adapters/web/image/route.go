package image

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	an "github.com/lejeunel/go-image-annotator/adapters/web/annotator"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	rt "github.com/lejeunel/go-image-annotator/routes"
)

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(mws...)

		r.Get(rt.ImagesUrl, s.List)
		r.Get(ImageRow, s.TableRow)
		r.Delete(ImageRow, s.Delete)

		r.Get(ingestPanelUrl, s.IngestionPanel)
		r.Post(archiveIngestUrl, s.IngestArchive)

		r.Post(rt.SliceUrl, s.Slice)
	})
}

type ImageURLFunc func(image im.Image) string

func BaseAnnotateImageURLFunc(image im.Image) string {
	base := MakeAnnotateImageURL(an.AnnotateImage, image.Id.String(), image.Collection.Name)
	withQuery := rt.AddQueryParams(base,
		rt.FilterQueryArgName, fmt.Sprintf("collection:\"%v\"", image.Collection.Name),
		rt.OrderingQueryArgName, "ingested_at")
	return withQuery.String()
}

func MakeAnnotateImageURLFunc(image im.Image, filters, ordering string) ImageURLFunc {
	return func(image im.Image) string {
		base := MakeAnnotateImageURL(an.AnnotateImage, image.Id.String(), image.Collection.Name)
		withQuery := rt.AddQueryParams(base,
			rt.FilterQueryArgName, filters,
			rt.OrderingQueryArgName, ordering)
		return withQuery.String()
	}
}
