package builders

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func BaseLibs() []Node {
	return []Node{
		Raw("<style>[x-cloak] { display: none !important; }</style>"),
		Script(
			Src("/static/htmx.js"),
			Defer(),
		),
		Script(
			Src("/static/alpine-focus.js"),
			Defer(),
		),
		Script(
			Src("/static/alpine.js"),
			Defer(),
		),
	}
}
