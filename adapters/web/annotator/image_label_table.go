package annotator

import (
	"fmt"
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	"github.com/lejeunel/go-image-annotator/modules/annotator/view"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ImageLabelRow struct {
	Label  string
	Id     string
	Author string
	Time   string
}

func (r ImageLabelRow) Render() Node {
	return Tr(Class("text-left"),
		Td(
			Div(Class("flex flex-col"),
				Div(Class(authorInfo), Text(r.Author)),
				Div(Class(authorInfo), Text(r.Time)),
			),
		),
		Td(Class("ps-2 py-1"),
			Text(r.Label),
		),
		Td(Div(
			Class("flex justify-end items-center pr-1"),
			Raw(fmt.Sprintf(`<a href="#" onclick="AnnotatorModule.remove('%v')"> %v </a>`, r.Id, ic.TrashIcon)),
		),
		))
}

type ImageLabelTable struct {
	Fields []string
	Rows   []ImageLabelRow
}

func (t *ImageLabelTable) AddImageLabel(l view.ImageLabel) {
	t.Rows = append(t.Rows, ImageLabelRow{Label: l.Label, Id: l.Id, Author: l.Author, Time: l.Time})
}

func (t *ImageLabelTable) Build() Node {
	return Div(Class("pb-2"),
		Div(Class("overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
			Table(Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
				TBody(Class("divide-y divide-outline dark:divide-outline-dark"),
					Tr(
						Td(Div(Class("text-left py-2 ps-2 pe-2 text-sm font-bold"), Text("Labels"))),
						Td(),
						Td(Class("align-middle"),
							Div(Class("flex items-center justify-end pr-1"),
								Raw(fmt.Sprintf(`<a href="#" onclick="Alpine.store('imageLabelModal').open()"> %v </a>`, ic.AddCircleIcon)),
							),
						),
					),
					Map(t.Rows, func(r ImageLabelRow) Node {
						return r.Render()
					}),
				)),
		))

}
