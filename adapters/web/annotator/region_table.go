package annotator

import (
	"bytes"
	"fmt"
	"github.com/lejeunel/go-image-annotator/modules/annotator/view"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"text/template"
)

type RegionTable struct {
	Rows            []RegionRow
	AvailableLabels []string
}

func (t *RegionTable) AddBox(b view.BoundingBox) {
	shortId := ShortenUUID(b.Id)
	tmpl := template.New("")
	template.Must(tmpl.ParseFS(templatesFiles,
		"templates/label_combobox.html"))

	var buf bytes.Buffer
	tmpl.ExecuteTemplate(&buf, "label_combobox",
		LabelSelector{Labels: t.AvailableLabels, SelectorIsOpen: false, Selected: &b.Label, AnnotationId: b.Id})
	t.Rows = append(t.Rows,
		RegionRow{Values: []Node{
			Div(Class("ps-1"),
				Raw(fmt.Sprintf(`<svg width="22" height="22" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="10" y="10" width="80" height="80" rx="10" ry="10" fill="%v" />
</svg>`, b.Color)),
			),
			Text(shortId),
			Raw(buf.String()),
			Div(
				Class("flex  justify-end items-center pr-1"),
				Raw(fmt.Sprintf(`<a href="#" onclick="AnnotatorModule.remove('%v')"> %v </a>`, b.Id, TrashIcon)),
			),
		}})

}

func (t *RegionTable) Build(title string) Node {
	return Div(Class("overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
		Table(Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
			RegionTableBody(title, t.Rows),
		))
}

type RegionRow struct {
	Values []Node
}

func (r RegionRow) Render() Node {
	return Tr(
		Map(r.Values, func(node Node) Node {
			return Td(node)
		}))
}

func RegionTableBody(title string, rows []RegionRow) Node {
	return TBody(Class("divide-y divide-outline dark:divide-outline-dark"),
		Td(Div(Class("text-left py-2 ps-2 pe-2 text-sm font-bold"), Text(title))),
		Map(rows, func(r RegionRow) Node {
			return r.Render()

		}),
	)
}
