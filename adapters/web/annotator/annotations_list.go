package annotator

import (
	"fmt"
	a "github.com/lejeunel/go-image-annotator-v2/application/annotator"
	html "github.com/lejeunel/go-image-annotator-v2/shared/html"
	. "maragu.dev/gomponents"
)

var TrashIcon = `
<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
</svg>
`

type AnnotationsListView struct{}

func (v *AnnotationsListView) Build(boxes []*a.BoundingBox) Node {
	table := html.PaginationTable{Fields: []string{"", "id", "label", "actions"}}
	for i, b := range boxes {
		shortId := b.Id[:8]
		table.Rows = append(table.Rows,
			html.PaginationTableRow{Values: []Node{
				Raw(fmt.Sprintf(`<svg width="22" height="22" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="10" y="10" width="80" height="80" rx="10" ry="10" fill="%v" />
</svg>`, a.Palette[i])),
				Text(shortId),
				Text(b.Label),
				Raw(fmt.Sprintf(`<a href="#" onclick="removeAnnotation('%v')"> %v </a>`, b.Id, TrashIcon))}})
	}
	return table.Build()
}
