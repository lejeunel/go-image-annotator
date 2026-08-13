package annotator

import (
	"fmt"
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	listmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/list"
	readmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/read"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type AddMetaPresenter struct {
	data   []m.MetaData
	Writer http.ResponseWriter
	GetURL b.RowURL
	htmx.ErrorPresenter
}

func NewAddMetaPresenter(w http.ResponseWriter, u b.RowURL) AddMetaPresenter {
	return AddMetaPresenter{
		Writer:         w,
		GetURL:         u,
		ErrorPresenter: htmx.NewErrorPresenter("adding meta-data", w),
	}
}

func (p *AddMetaPresenter) SuccessAddMetadata(r m.MetaData) {
	htmx.NotifySuccessPayload(p.Writer,
		"adding meta-data",
		fmt.Sprintf("Successfully added %v:%v", r.Key, r.Value))
}

func (p *AddMetaPresenter) SuccessListMetadata(r listmd.Response) {
	RenderMetaDataList(r.ImageId, r.Collection, r.MetaData, p.Writer)
}

type MetaRowPresenter struct {
	writer http.ResponseWriter
	htmx.ErrorPresenter
}

func NewMetaRowPresenter(w http.ResponseWriter) MetaRowPresenter {
	return MetaRowPresenter{w, htmx.NewErrorPresenter("viewing meta-data", w)}
}

func (p MetaRowPresenter) SuccessReadMetadata(r readmd.Response) {
	MakeRow(r.ImageId, r.Collection, r.Data).Render(p.writer)
}

type MetaDeletePresenter struct {
	writer http.ResponseWriter
	b.RowURL
	htmx.ErrorPresenter
}

func NewMetaDeletePresenter(w http.ResponseWriter, u b.RowURL) MetaDeletePresenter {
	return MetaDeletePresenter{w, u, htmx.NewErrorPresenter("viewing meta-data", w)}
}

func (p MetaDeletePresenter) SuccessDeleteMetadata(key string) {
	htmx.NotifySuccessPayload(p.writer,
		"deleting meta-data item",
		fmt.Sprintf("Successfully deleted key %v", key))
}

func BuildMetaDataList(imageId, collection string, data []m.MetaData) Node {
	table := tb.NewTableBuilder(MetaTableFields, tb.WithSimplePlaceHolder())
	for _, d := range data {
		table.AddRow(MakeRow(imageId, collection, d))
	}

	buttons := Div(
		Class("flex gap-1"),
		Group([]Node{
			MakeAddMetaHTMXButton(
				ic.Text,
				"Add text",
				MetaTypeString,
				imageId,
				collection,
			),
			MakeAddMetaHTMXButton(
				ic.Hash,
				"Add integer",
				MetaTypeInt,
				imageId,
				collection,
			),
			MakeAddMetaHTMXButton(
				ic.Percent,
				"Add float",
				MetaTypeFloat,
				imageId,
				collection,
			),
			MakeAddMetaHTMXButton(
				ic.Calendar,
				"Add date/time",
				MetaTypeDateTime,
				imageId,
				collection,
			),
			MakeAddMetaHTMXButton(
				ic.Flag,
				"Add boolean",
				MetaTypeBoolean,
				imageId,
				collection,
			),
		}))
	return Div(
		ID(MetaDivId),
		cmp.Separator,
		Div(Class("flex items-center mt-2 mb-2"),
			Div(Class("flex-1 text-lg"), Text("Custom Meta-data")),
			buttons,
		),
		Div(ID(MetaFormId)),
		table.Build())
}

func RenderMetaDataList(imageId, collection string, data []m.MetaData, w io.Writer) {
	BuildMetaDataList(imageId, collection, data).Render(w)
}

func MakeRow(imageId, collection string, data m.MetaData) tb.Row {
	url := b.NewRowURL(MetaRowUrl)
	url.Set(MetaImageIdArg, imageId)
	url.Set(MetaCollectionArg, collection)
	url.Set(MetaKeyArg, data.Key)

	actions := b.NewActionsPanelBuilder()
	actions.SetConfirmDelete(url.SetMode(b.ModeConfirmDelete).Url)
	r := tb.NewRow()
	r.AddCell(tb.NewCell(Text(data.Key)))
	r.AddCell(tb.NewCell(Text(fmt.Sprint(data.Value))))
	r.AddCell(tb.NewCell(actions.Build()))
	return r
}
