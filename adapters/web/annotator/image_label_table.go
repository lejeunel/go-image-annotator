package annotator

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/app/annotator/view"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ImageLabelRow struct {
	Label string
	Id    string
}

func (r ImageLabelRow) Render() Node {
	return Tr(Class("text-left"), Td(Class("ps-2 py-1"), Text(r.Label)), Td(
		Div(
			Class("pe-1 text-right"),
			Raw(fmt.Sprintf(`<a href="#" onclick="AnnotatorModule.remove('%v')"> %v </a>`, r.Id, TrashIcon)),
		),
	))
}

type ImageLabelTable struct {
	Fields []string
	Rows   []ImageLabelRow
}

func (t *ImageLabelTable) AddImageLabel(l view.ImageLabel) {
	t.Rows = append(t.Rows, ImageLabelRow{l.Label, l.Id})
}

func (t *ImageLabelTable) Build() Node {
	return Div(Class("pb-2"),
		Div(Class("overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
			Table(Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
				TBody(Class("divide-y divide-outline dark:divide-outline-dark"),
					Tr(
						Td(Div(Class("text-left py-2 ps-2 pe-2 text-sm font-bold"), Text("Labels"))),
						Td(Class("text-right"),
							Raw(fmt.Sprintf(`<a href="#" onclick="Alpine.store('imageLabelModal').open()"> %v </a>`, AddIcon)),
						),
					),
					Map(t.Rows, func(r ImageLabelRow) Node {
						return r.Render()
					}),
				)),
		))

}
