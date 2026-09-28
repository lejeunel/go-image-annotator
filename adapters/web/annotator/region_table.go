package annotator

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	"github.com/lejeunel/go-image-annotator/modules/annotator/view"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type LabelSelector struct {
	Labels       []string
	Selected     *string
	AnnotationId string
}

type RegionKind int

const (
	RegionBox RegionKind = iota
	RegionPolygon
)

var authorInfo = "text-xs italic text-gray-500 dark:gray-500 ml-1"

type RegionTable struct {
	Rows            []RegionRow
	AvailableLabels []string
}

func (t *RegionTable) addRow(author, time, id, label, color string, regionKind RegionKind) {
	var regionIcon string
	switch regionKind {
	case RegionBox:
		regionIcon = ic.MakeColoredRectangleIcon(color)
	case RegionPolygon:
		regionIcon = ic.MakeColoredHexagonIcon(color)
	}
	t.Rows = append(
		t.Rows,
		RegionRow{Author: author, Time: time, Id: id, Label: label, Color: color, Icon: regionIcon},
	)
}

func (t *RegionTable) SortRows() {
	sort.Slice(t.Rows, func(i int, j int) bool {
		return t.Rows[i].Time < t.Rows[j].Time
	})
}

func (t *RegionTable) AddPolygon(p view.Polygon) {
	t.addRow(p.Author, p.Time, p.Id, p.Label, p.Color, RegionPolygon)
}

func (t *RegionTable) AddBox(b view.BoundingBox) {
	t.addRow(b.Author, b.Time, b.Id, b.Label, b.Color, RegionBox)
}

func (t *RegionTable) Build(title string) Node {
	t.SortRows()
	return Div(
		Class(
			"w-full rounded-radius border border-outline dark:border-outline-dark",
		),
		Table(Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
			RegionTableBody(title, t.Rows),
		),
	)
}

type RegionRow struct {
	Author string
	Label  string
	Time   string
	Id     string
	Icon   string
	Color  string
	Date   string
}

func (r RegionRow) Render(w io.Writer) {
	Div(
		Td(Class("ps-1 py-2"),
			Div(Class("flex flex-col"),
				Div(Class(authorInfo), Text(r.Author)),
				Div(Class(authorInfo), Text(r.Time)),
			),
		),
		Td(Class("ps-1 py-2"),
			Raw(r.Icon),
		),
		Td(Text(r.Label)),
		Td(
			Class("flex  justify-end items-center pr-1 gap-1 ps-1 py-2"),
			cmp.MakeIconizedButton(ic.Edit, "edit",
				Attr(fmt.Sprintf(
					`onclick="
						Annotator.setAnnotationId('%v');
						Annotator.editLabelMode();
						LabelPicker.open();"`,
					r.Id))),
			cmp.MakeIconizedButton(
				ic.Trash,
				"delete",
				Attr(fmt.Sprintf("onclick=\"Annotator.remove('%v')\"", r.Id)),
			),
		)).Render(w)
}

func RegionTableBody(title string, rows []RegionRow) Node {
	return TBody(Class("divide-y divide-outline dark:divide-outline-dark"),
		Td(Div(Class("text-left py-2 ps-2 pe-2 text-sm font-bold"), Text(title))),
		Map(rows, func(r RegionRow) Node {
			var buf bytes.Buffer
			r.Render(&buf)
			return Tr(Raw(buf.String()))
		}),
	)
}
