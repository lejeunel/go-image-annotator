package annotator

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var TrashIcon = `
<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
</svg>
`

var Badge = `<span class="w-fit inline-flex mx-1 my-1 overflow-hidden rounded-radius border border-secondary bg-surface text-xs font-medium text-secondary dark:border-secondary-dark dark:bg-surface-dark dark:text-secondary-dark">
    <span class="flex items-center gap-1 bg-secondary/10 px-2 py-1 dark:bg-secondary-dark/10">
		<a href="#" onclick="%v">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" aria-hidden="true" stroke="currentColor" fill="none" stroke-width="1.4" class="size-4">
				<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
			</svg>
		</a>
		%v
    </span>
</span>`

func MakeImageLabelBadge(label string, id string) Node {
	return Raw(fmt.Sprintf(Badge, fmt.Sprintf("AnnotatorModule.remove('%v')", id), label))
}

func ShortenUUID(id string) string {
	return id[:8]
}

type AnnotationTable struct {
	Fields []string
	Rows   []AnnotationRow
}

func (t *AnnotationTable) Build() Node {
	return Div(Class("overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
		Table(Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
			TableHeader(t.Fields),
			TableBody(t.Rows),
		))
}

type AnnotationRow struct {
	Values []Node
}

func (r AnnotationRow) Render() Node {
	return Tr(
		Class("even:bg-primary/5 dark:even:bg-primary-dark/10"),
		Map(r.Values, func(node Node) Node {
			return Td(Class("p-2"),
				node)
		}))

}

func TableHeader(fields []string) Node {
	return THead(Tr(Class("border-b border-outline bg-surface-alt text-sm text-on-surface-strong dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark-strong"),
		Map(fields, func(f string) Node {
			return Th(Scope("col"), Class("p-2"), Text(f))
		})))
}

func TableBody(rows []AnnotationRow) Node {
	return TBody(Class("divide-y divide-outline dark:divide-outline-dark"),
		Map(rows, func(r AnnotationRow) Node {
			return r.Render()

		}),
	)
}

type AnnotationsListView struct{}

func (v *AnnotationsListView) makeBoxList(boxes []*view.BoundingBox) Node {
	bboxTable := AnnotationTable{Fields: []string{"", "id", "label", "actions"}}
	for i, b := range boxes {
		shortId := ShortenUUID(b.Id)
		bboxTable.Rows = append(bboxTable.Rows,
			AnnotationRow{Values: []Node{
				Raw(fmt.Sprintf(`<svg width="22" height="22" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="10" y="10" width="80" height="80" rx="10" ry="10" fill="%v" />
</svg>`, view.Palette[i])),
				Text(shortId),
				Text(b.Label),
				Raw(fmt.Sprintf(`<a href="#" onclick="AnnotatorModule.remove('%v')"> %v </a>`, b.Id, TrashIcon))}})
	}
	return bboxTable.Build()
}

func (v *AnnotationsListView) makeImageLabelList(imageLabels []*view.ImageLabel) Node {
	badges := []Node{}
	for _, l := range imageLabels {
		badges = append(badges, MakeImageLabelBadge(l.Label, l.Id))
	}
	return Div(Class("w-64 py-4"), Div(badges...))
}

func (v *AnnotationsListView) Build(boxes []*view.BoundingBox, imageLabels []*view.ImageLabel) Node {
	bboxTable := v.makeBoxList(boxes)
	imageLabelsTable := v.makeImageLabelList(imageLabels)

	fullTable := Div(
		imageLabelsTable,
		bboxTable)
	return fullTable
}
