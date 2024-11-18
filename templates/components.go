package templates

import (
	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func page(title string, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []Node{
			Script(Src("https://cdn.tailwindcss.com?plugins=typography")),
			Script(Defer(), Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js")),
		},
		Body: []Node{Class("bg-gradient-to-b from-white to-indigo-100 bg-no-repeat"),
			Div(Class("min-h-screen flex flex-col justify-between"),
				NavBar(),
				Div(Class("grow"),
					container(true,
						Div(Class("prose prose-lg prose-indigo"),
							Group(children),
						),
					),
				),
				footer(),
			),
		},
	})
}

func NavBar() Node {
	return Nav(
		x.Data(`{navigationMenuOpen: false,
        navigationMenu: '',
        navigationMenuCloseDelay: 200,
        navigationMenuCloseTimeout: null,
        navigationMenuLeave() {
            let that = this;
            this.navigationMenuCloseTimeout = setTimeout(() => {
                that.navigationMenuClose();
            }, this.navigationMenuCloseDelay);
        },
        navigationMenuReposition(navElement) {
            this.navigationMenuClearCloseTimeout();
            this.$refs.navigationDropdown.style.left = navElement.offsetLeft + 'px';
            this.$refs.navigationDropdown.style.marginLeft = (navElement.offsetWidth/2) + 'px';
        },
        navigationMenuClearCloseTimeout(){
            clearTimeout(this.navigationMenuCloseTimeout);
        },
        navigationMenuClose(){
            this.navigationMenuOpen = false;
            this.navigationMenu = '';
        }
    }`), Class("relative z-10 w-auto"),
		Div(Class("relative"),
			Ul(Class("flex items-center justify-center flex-1 p-1 space-x-1 list-none border rounded-md text-neutral-700 group border-neutral-200/80"),
				Li(A(Href("#_"),
					Class("inline-flex items-center justify-center h-10 px-4 py-2 text-sm font-medium transition-colors rounded-md hover:text-neutral-900 focus:outline-none disabled:opacity-50 disabled:pointer-events-none bg-background hover:bg-neutral-100 group w-max"),
					Text("Documentation"))))))
}

func headerLink(href, text string) Node {
	return A(Class("hover:text-indigo-300"), Href(href), Text(text))
}

func container(padY bool, children ...Node) Node {
	return Div(
		Classes{
			"max-w-7xl mx-auto":     true,
			"px-4 md:px-8 lg:px-16": true,
			"py-4 md:py-8":          padY,
		},
		Group(children),
	)
}

func footer() Node {
	return Div(Class("bg-gray-900 text-white shadow text-center h-16 flex items-center justify-center"),
		A(Href("https://www.gomponents.com"), Text("gomponents")),
	)
}
