package web

import (
	"fmt"
	"net/http"
	"net/url"

	aw "github.com/lejeunel/go-image-annotator-v2/adapters/web/annotator"
	an "github.com/lejeunel/go-image-annotator-v2/application/annotator"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	"github.com/lejeunel/go-image-annotator-v2/shared/html"
)

type Request struct {
	ImageId    im.ImageId
	Collection string
}

func ParseURL(u *url.URL) (*Request, error) {
	baseErr := "parsing url"
	req := Request{}
	imageIdStr := u.Query().Get("id")
	if imageIdStr == "" {
		return nil, fmt.Errorf("%v: extracting id: %w", baseErr, e.ErrURLParsing)
	}
	imageId, err := im.NewImageIdFromString(imageIdStr)
	if err != nil {
		return nil, fmt.Errorf("%v: validating id (%v): %w", baseErr, imageIdStr, e.ErrValidation)

	}
	req.ImageId = imageId

	collection := u.Query().Get("collection")
	if collection == "" {
		return nil, fmt.Errorf("%v: collection (%v): %w", baseErr, collection, e.ErrURLParsing)
	}
	req.Collection = collection
	return &req, nil
}

func (s *Server) ViewImage(w http.ResponseWriter, r *http.Request) {

	req, err := ParseURL(r.URL)
	if err != nil {
		html.NewPageBuilder().SetError(err).Render(w)
		return
	}
	view := aw.NewAnnotationView()
	presenter := an.NewAnnotatorPresenter(view)
	s.annotator.Start(req.ImageId, req.Collection, *presenter)
	view.Render(w)
}
