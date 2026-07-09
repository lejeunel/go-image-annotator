package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakeCard(node Node) Node {
	return Article(
		Class("group rounded-radius flex max-w-md flex-col border border-outline bg-surface-alt p-2 text-on-surface dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark mt-6 mb-6"),
		P(Class("text-pretty text-sm"), node),
	)
}
