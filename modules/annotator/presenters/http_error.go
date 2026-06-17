package presenters

import (
	"io"

	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	updpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-polygon"
	del "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

type HTTPErrorPresenter struct {
	Err error
}

func NewHTTPErrorPresenter() HTTPErrorPresenter {
	return HTTPErrorPresenter{}
}

func (p HTTPErrorPresenter) SuccessAddLabel(r addlbl.Response)       {}
func (p HTTPErrorPresenter) SuccessAddBox(r addbox.Response)         {}
func (p HTTPErrorPresenter) SuccessAddPolygon(r addpoly.Response)    {}
func (p HTTPErrorPresenter) SuccessUpdatePolygon(r updpoly.Response) {}
func (p HTTPErrorPresenter) SuccessUpdateBox(r updbox.Response)      {}
func (p HTTPErrorPresenter) SuccessUpdateLabel(r updlbl.Response)    {}
func (p HTTPErrorPresenter) SuccessDeleteAnnotation(r del.Response)  {}
func (p *HTTPErrorPresenter) Error(err error) {
	p.Err = err
}
func (p HTTPErrorPresenter) Write(w io.Writer) {
	// TODO write http error here
}
