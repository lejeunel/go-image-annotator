package components

import (
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CopyURLButton() Node {
	return Button(
		Type("button"),
		Title("Copy URL to clipboard"),
		Aria("label", "Copy URL to clipboard"),
		Class(
			"flex cursor-pointer [&_svg]:size-5 text-on-surface hover:opacity-75 dark:text-on-surface-dark",
		),
		Attr("x-data", "{ copied: false }"),
		Attr(
			"@click",
			"navigator.clipboard.writeText(window.location.href).then(() => { copied = true; setTimeout(() => copied = false, 2000) })",
		),
		Span(Attr("x-show", "!copied"), Raw(ic.Clipboard)),
		Span(Attr("x-show", "copied"), Attr("x-cloak", ""), Class("text-success"), Raw(ic.Check)),
	)
}
