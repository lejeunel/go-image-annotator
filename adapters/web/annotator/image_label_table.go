package annotator

import (
	"fmt"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
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
		Td(Class("ps-2 py-3 min-w-0 break-words"),
			Text(r.Label),
		),
		Td(Div(
			Class("flex justify-end items-center pr-1 gap-1"),
			cmp.MakeIconizedButton(ic.Edit, "edit",
				Attr(fmt.Sprintf(
					`onclick="
						Annotator.setAnnotationId('%v');
						Annotator.editLabelMode();
						LabelPicker.open();"`,
					r.Id))),
			cmp.MakeIconizedButton(ic.Trash, "delete", Attr(fmt.Sprintf("onclick=\"Annotator.remove('%v')\"", r.Id))),
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
	return Div(
		Class("pb-2"),
		Div(
			Class(
				"overflow-visible w-full rounded-radius border border-outline dark:border-outline-dark",
			),
			Table(Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
				TBody(Class("divide-y divide-outline dark:divide-outline-dark"),
					Tr(
						Td(
							Div(
								Class("text-left py-2 ps-2 pe-2 text-sm font-bold"),
								Text("Labels"),
							),
						),
						Td(),
						Td(Class("align-middle"),
							Div(
								Class("flex items-center justify-end pr-1 py-1"),
								cmp.MakeIconizedButton(ic.Add, "Add image label",
									Attr(`onclick="
											Annotator.imageLabelMode();
											Annotation.LabelPicker.open();
								"`),
								),
							),
						),
					),
					Map(t.Rows, func(r ImageLabelRow) Node {
						return r.Render()
					}),
				)),
		),
	)
}
