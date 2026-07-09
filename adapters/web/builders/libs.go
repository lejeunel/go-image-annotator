package builders

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func BaseJSDeps() []Node {
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
		Script(Raw(`
			function notify(variant, title, message, extra = {}) {
				window.dispatchEvent(new CustomEvent("notify", {
					detail: {
						variant: variant,
						title: title,
						message: message,
						...extra,
					},
				}));
			}
`)),
	}
}

func BaseBodyExtra() string {
	notificationArea, err := componentsFiles.ReadFile("components/notifications.html")
	if err != nil {
		panic(err)
	}
	return string(notificationArea)
}
