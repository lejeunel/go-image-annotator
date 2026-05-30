package annotator

import (
	"bytes"
	"fmt"
	"github.com/lejeunel/go-image-annotator/app/annotator/view"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"text/template"
)

var TrashIcon = `
<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
</svg>
`

var BadgeIcon = `<span class="w-fit inline-flex mx-1 my-1 overflow-hidden rounded-radius border border-secondary bg-surface text-xs font-medium text-secondary dark:border-secondary-dark dark:bg-surface-dark dark:text-secondary-dark">
    <span class="flex items-center gap-1 bg-secondary/10 px-2 py-1 dark:bg-secondary-dark/10">
		<a href="#" onclick="%v">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" aria-hidden="true" stroke="currentColor" fill="none" stroke-width="1.4" class="size-4">
				<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
			</svg>
		</a>
		%v
    </span>
</span>`

var AddIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <line x1="12" y1="4" x2="12" y2="20"/>
  <line x1="4" y1="12" x2="20" y2="12"/>
</svg>`

func MakeImageLabelBadge(label string, id string) Node {
	return Raw(fmt.Sprintf(BadgeIcon, fmt.Sprintf("AnnotatorModule.remove('%v')", id), label))
}

func ShortenUUID(id string) string {
	return id[:8]
}

type AnnotationTable struct {
	Fields []string
	Rows   []AnnotationRow
}

func (t *AnnotationTable) Build(title string) Node {
	return Div(Class("overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
		Table(Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
			TableBody(title, t.Rows),
		))
}

type AnnotationRow struct {
	Values []Node
}

func (r AnnotationRow) Render() Node {
	return Tr(
		Map(r.Values, func(node Node) Node {
			return Td(Class("p-1"),
				node)
		}))

}

func TableHeader(fields []string) Node {
	return THead(Tr(Class("border-b border-outline bg-surface-alt text-sm text-on-surface-strong dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark-strong"),
		Map(fields, func(f string) Node {
			return Th(Scope("col"), Class("p-2"), Text(f))
		})))
}

func TableBody(title string, rows []AnnotationRow) Node {
	return TBody(Class("divide-y divide-outline dark:divide-outline-dark"),
		Td(Div(Class("text-left py-2 ps-2 pe-2 text-sm font-bold"), Text(title))),
		Map(rows, func(r AnnotationRow) Node {
			return r.Render()

		}),
	)
}

type LabelSelector struct {
	Labels         []string
	SelectorIsOpen bool
	Selected       *string
	AnnotationId   string
}
type AnnotationsListView struct{}

func (v *AnnotationsListView) makeBoxList(boxes []view.BoundingBox, availableLabels []string) Node {
	bboxTable := AnnotationTable{Fields: []string{"", "id", "label", "actions"}}
	for _, b := range boxes {
		shortId := ShortenUUID(b.Id)
		t := template.New("")
		template.Must(t.ParseFS(templatesFiles,
			"templates/label_combobox.html"))

		var buf bytes.Buffer
		t.ExecuteTemplate(&buf, "label_combobox",
			LabelSelector{Labels: availableLabels, SelectorIsOpen: false, Selected: &b.Label, AnnotationId: b.Id})
		bboxTable.Rows = append(bboxTable.Rows,
			AnnotationRow{Values: []Node{
				Raw(fmt.Sprintf(`<svg width="22" height="22" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="10" y="10" width="80" height="80" rx="10" ry="10" fill="%v" />
</svg>`, b.Color)),
				Text(shortId),
				Raw(buf.String()),
				Raw(fmt.Sprintf(`<a href="#" onclick="AnnotatorModule.remove('%v')"> %v </a>`, b.Id, TrashIcon))}})
	}
	return bboxTable.Build("Regions")
}

func (v *AnnotationsListView) makeImageLabelList(imageLabels []view.ImageLabel) Node {
	badges := []Node{}
	for _, l := range imageLabels {
		badges = append(badges, MakeImageLabelBadge(l.Label, l.Id))
	}
	return Div(Class("w-64 py-4"), Div(badges...))
}

func (v *AnnotationsListView) Build(boxes []view.BoundingBox, imageLabels []view.ImageLabel, availableLabels []string) Node {
	bboxTable := v.makeBoxList(boxes, availableLabels)

	var imageLabelsTable Node
	if len(imageLabels) > 0 {
		imageLabelsTable = v.makeImageLabelList(imageLabels)
	}
	imageLabelsPicker := Div(Class("pb-2"),
		Div(Class("rounded-radius border border-outline dark:border-outline-dark"),
			Table(
				Tr(
					Td(Div(Class("text-left py-2 ps-2 pe-4 text-sm font-bold"), Text("Labels"))),
					Td(Class("text-right"),
						Raw(fmt.Sprintf(`<a href="#" onclick=""> %v </a>`, AddIcon)),
					),
				),
			)))

	fullTable := Div(
		imageLabelsPicker,
		imageLabelsTable,
		bboxTable)
	return fullTable
}
