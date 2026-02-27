package images

import (
	"bufio"
	"bytes"
	"context"
	clc "datahub/domain/collections"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	g "datahub/generic"
	"fmt"
	"io"
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	"net/http"
	"net/url"
)

type ImagesListViewer struct {
	Title                string
	GenericViewer        *g.GenericViewer
	ImageSetDescriber    ImageSetDescriber
	ImageService         *Service
	PageSize             int
	PaginatorBaseURL     *url.URL
	PaginationWidgetSize int
}

func NewImagesListViewer(title string, gv *g.GenericViewer, d ImageSetDescriber,
	s *Service, pagesize int, paginatorWidgetSize int, paginatorBaseURL string) *ImagesListViewer {

	baseURL, _ := url.Parse(paginatorBaseURL)
	return &ImagesListViewer{
		Title:                title,
		GenericViewer:        gv,
		ImageSetDescriber:    d,
		ImageService:         s,
		PageSize:             pagesize,
		PaginatorBaseURL:     baseURL,
		PaginationWidgetSize: paginatorWidgetSize,
	}

}

func (v *ImagesListViewer) RenderPaginatedList(ctx context.Context, images []Image, p *g.PaginationMeta,
	originEntity string, originId string) (gp.Node, error) {
	paginationWidget, err := g.MakePaginationViewer(*p,
		*v.PaginatorBaseURL,
		int(p.PageSize),
		v.PaginationWidgetSize)
	if err != nil {
		return nil, err
	}

	var rows []g.TableRow
	for _, im := range images {

		isAnnotated := "false"

		if len(im.BoundingBoxes) > 0 {
			isAnnotated = "true"
		}

		rows = append(rows, g.TableRow{[]gp.Node{
			gh.A(gh.Href(fmt.Sprintf("/image?collection_id=%v&image_id=%v&origin_entity=%v&origin_id=%v&ordering=%v&descending=%v",
				im.Collection.Id, im.Id.String(), originEntity, originId, "captured_at", "false",
			)),
				gh.Class(g.UrlClass),
				gh.Span(gp.Text(im.Id.String()))),
			gp.Raw(im.Collection.Name),
			gp.Raw(isAnnotated),
			gp.Raw(im.GetSiteName()),
			gp.Raw(im.GetCameraName()),
			gp.Raw(im.CapturedAt.Format("2006-01-02 15:04:05.000")),
			gp.Raw(im.CreatedAt.Format("2006-01-02 15:04:05.000")),
		}})
	}

	imageTableList := g.MyTable([]string{"id", "collection", "annotated?", "site", "camera", "captured_at", "created_at"},
		rows)

	return gp.Group([]gp.Node{
		*paginationWidget,
		imageTableList,
		*paginationWidget}), nil

}

func (v *ImagesListViewer) Render(ctx context.Context,
	filters FilterArgs, ordering OrderingArgs,
	pagination g.PaginationParams,
	originEntity string, originId string,
	w io.Writer) {

	images, paginationMeta, err := v.ImageService.List(ctx, filters, ordering, pagination,
		FetchMetaOnly)
	if err != nil {
		v.GenericViewer.OopsPage(ctx, err.Error()).Render(w)
		return
	}

	paginatedView, err := v.RenderPaginatedList(ctx, images, paginationMeta, originEntity, originId)
	if err != nil {
		v.GenericViewer.OopsPage(ctx, err.Error()).Render(w)
		return
	}

	var paginatedViewWithTitle gp.Group
	if len(images) > 0 {
		paginatedViewWithTitle = gp.Group{gh.H2(gh.Class("text-3xl font-bold text-gray-900 py-2"), gp.Text("Images")),
			paginatedView}
	} else {
		paginatedViewWithTitle = gp.Group{gh.H2(gh.Class("text-3xl font-bold text-gray-900 py-2"), gp.Text("Images")),
			gp.Text("Oh so empty (◡︵◡)")}
	}

	var descriptionBuffer bytes.Buffer
	descriptionWriter := bufio.NewWriter(&descriptionBuffer)
	if err := v.ImageSetDescriber.Describe(ctx, descriptionWriter); err != nil {
		v.GenericViewer.OopsPage(ctx, err.Error()).Render(w)
		return
	}
	if err := descriptionWriter.Flush(); err != nil {
		v.GenericViewer.OopsPage(ctx, err.Error()).Render(w)
		return
	}

	var imageListBuffer bytes.Buffer
	imageListWriter := bufio.NewWriter(&imageListBuffer)
	list := gh.Div(gh.Class("px-2 md:px-4 lg:px8 py-2 md:py-4"), paginatedViewWithTitle)
	list.Render(imageListWriter)
	if err := imageListWriter.Flush(); err != nil {
		v.GenericViewer.OopsPage(ctx, err.Error()).Render(w)
		return
	}

	v.GenericViewer.BasePage(ctx, v.Title, g.Head(), gp.Group{gp.Raw(descriptionBuffer.String()),
		gp.Raw(imageListBuffer.String())}).Render(w)
}

func NewImageListOfCollectionHandler(gv *g.GenericViewer, imageService *Service, collectionService *clc.Service, pageSize int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		collectionId, err := clc.NewCollectionIdFromString(r.PathValue("id"))
		if err != nil {
			gv.OopsPage(r.Context(), err.Error()).Render(w)
			return
		}

		filters := FilterArgs{CollectionId: collectionId}
		describer := clc.CollectionDescriber{CollectionService: collectionService, Id: *collectionId}
		pagination, err := g.NewPaginationParamsFromURL(*r.URL, pageSize)
		if err != nil {
			gv.OopsPage(r.Context(), err.Error()).Render(w)
			return
		}
		viewer := NewImagesListViewer("Collection", gv, describer, imageService, pageSize, 10,
			fmt.Sprintf("/collection/%v", collectionId))
		viewer.Render(r.Context(), filters, *NewImageDefaultOrderingArgs(), *pagination, "collection", collectionId.String(), w)
	})
}

func NewImageListOfCameraHandler(gv *g.GenericViewer, imageService *Service, locationService *loc.Service, pageSize int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cameraId, err := loc.NewCameraIdFromString(r.PathValue("id"))
		if err != nil {
			gv.OopsPage(r.Context(), err.Error()).Render(w)
			return
		}
		filters := FilterArgs{CameraId: cameraId}
		describer := loc.CameraDescriber{LocationService: locationService, Id: *cameraId}
		pagination, err := g.NewPaginationParamsFromURL(*r.URL, pageSize)
		if err != nil {
			gv.OopsPage(r.Context(), err.Error()).Render(w)
			return
		}
		viewer := NewImagesListViewer("Camera", gv, describer, imageService, pageSize, 10,
			fmt.Sprintf("/camera/%v", cameraId))
		viewer.Render(r.Context(), filters, *NewImageDefaultOrderingArgs(), *pagination, "camera", cameraId.String(), w)
	})
}

func NewImageListOfLabelHandler(gv *g.GenericViewer, imageService *Service, labelService *lbl.Service, pageSize int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		labelId, err := lbl.NewLabelIdFromString(r.PathValue("id"))
		if err != nil {
			gv.OopsPage(r.Context(), err.Error()).Render(w)
			return
		}
		filters := FilterArgs{LabelId: labelId}
		describer := lbl.LabelDescriber{LabelService: labelService, Id: *labelId}
		pagination, err := g.NewPaginationParamsFromURL(*r.URL, pageSize)
		if err != nil {
			gv.OopsPage(r.Context(), err.Error()).Render(w)
			return
		}
		viewer := NewImagesListViewer("Label", gv, describer, imageService, pageSize, 10,
			fmt.Sprintf("/label/%v", labelId))
		viewer.Render(r.Context(), filters, *NewImageDefaultOrderingArgs(), *pagination, "label", labelId.String(), w)
	})
}
