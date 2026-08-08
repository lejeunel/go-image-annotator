package components

import (
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakeClass(active bool) string {
	if active {
		return "flex items-center rounded-radius p-1 text-on-surface hover:text-primary dark:text-on-surface-dark dark:hover:text-primary-dark"
	}
	return "flex items-center rounded-radius p-1 text-gray-500 dark:text-gray-500 "
}

var leftArrow = `
					<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" class="size-6">
						<path fill-rule="evenodd" d="M11.78 5.22a.75.75 0 0 1 0 1.06L8.06 10l3.72 3.72a.75.75 0 1 1-1.06 1.06l-4.25-4.25a.75.75 0 0 1 0-1.06l4.25-4.25a.75.75 0 0 1 1.06 0Z" clip-rule="evenodd" />
					</svg>

				`

var rightArrow = `
					<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" class="size-6">
						<path fill-rule="evenodd" d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" />
					</svg>

				`

func MakeNavigationButton(
	url string,
	active bool,
	direction scr.ScrollingDirection,
	text string,
) Node {
	if direction == scr.ScrollNext {
		return A(Href(url), Class(MakeClass(active)),
			Text(text),
			Raw(rightArrow),
		)
	}
	return A(Href(url), Class(MakeClass(active)),
		Raw(leftArrow),
		Text(text),
	)
}

func MakeIconizedButton(icon, tooltip string, attrs ...Node) Node {
	return Div(
		Class("relative w-fit"),
		Button(
			Type("button"),
			Class(
				"peer rounded-radius bg-surface-alt border border-surface-alt px-1 py-1 font-medium tracking-wide text-on-surface focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:bg-surface-dark-alt dark:border-surface-dark-alt dark:text-on-surface-dark dark:focus-visible:outline-primary-dark cursor-pointer",
			),
			Group(attrs),
			Raw(icon),
		),
		Div(
			Class(
				"absolute -top-9 left-1/2 -translate-x-1/2 z-10 whitespace-nowrap rounded-sm bg-surface-dark px-2 py-1 text-center text-sm text-on-surface-dark-strong opacity-0 transition-all ease-out peer-hover:opacity-100 peer-focus:opacity-100 pointer-events-none dark:bg-surface dark:text-on-surface-strong",
			),
			Role("tooltip"),
			Text(tooltip),
		))
}
