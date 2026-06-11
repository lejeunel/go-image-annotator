package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakeTextLink(url, text string) Node {
	return A(
		Href(url),
		Class("font-medium text-primary underline-offset-2 hover:underline focus:underline focus:outline-hidden dark:text-primary-dark"),
		Text(text),
	)
}
