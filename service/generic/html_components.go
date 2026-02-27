package generic

import (
	"context"
	au "datahub/app/authorizer"
	"fmt"
	x "github.com/glsubri/gomponents-alpine"
	gp "maragu.dev/gomponents"
	gc "maragu.dev/gomponents/components"
	gh "maragu.dev/gomponents/html"
)

var navBarStyle = "inline-flex items-center justify-center h-10 px-4 py-2 text-sm font-medium transition-colors rounded-md hover:text-neutral-900 focus:outline-none disabled:opacity-50 disabled:pointer-events-none bg-background hover:bg-neutral-100 group w-max"
var navBarTitleStyle = "inline-flex items-center justify-center h-10 px-4 py-2 text-lg font-bold transition-colors rounded-md hover:text-neutral-900 focus:outline-none disabled:opacity-50 disabled:pointer-events-none bg-background hover:bg-neutral-100 group w-max"

func Head() gp.Node {
	return gp.Group([]gp.Node{
		gh.Link(gh.Href("/static/styles.css"), gh.Rel("stylesheet")),
		gh.Script(gh.Defer(), gh.Src("/static/alpine.js")),
		gh.Script(gh.Defer(), gh.Src("/static/htmx.js")),
	})
}

func BasePage(ctx context.Context, title string, head, body gp.Node,
	identityProvider au.IdentityProvider, signOutURL string) gp.Node {
	return HTMLPage("DataHub / "+title,
		head,
		gh.Div(gh.Class("min-h-screen flex flex-col justify-between"),
			NavBar(ctx, identityProvider, signOutURL),
			gh.Div(gh.Class("grow w-full px-2 md:px-4 lg:px-8 py-2 md:py-4"),
				body,
			),
		))
}

func HTMLPage(title string, head, body gp.Node) gp.Node {
	return gc.HTML5(gc.HTML5Props{
		Title:    title,
		Language: "en",
		Head:     []gp.Node{head},
		Body:     []gp.Node{body},
	})
}

func Prose(nodes ...gp.Node) gp.Node {
	return gh.Div(gh.Class("prose prose-lg prose-indigo max-w-none"),
		gp.Group(nodes),
	)
}

func NavBar(ctx context.Context, identityProvider au.IdentityProvider, signOutURL string) gp.Node {
	var username string
	username, err := identityProvider.Username(ctx)
	if err != nil {
		username = "<err-username>"

	}

	var email string
	email, err = identityProvider.Email(ctx)
	if err != nil {
		email = "<err-email>"

	}

	var entitlements []string
	entitlements, err = identityProvider.Entitlements(ctx)
	if err != nil {
		entitlements = []string{"<err-entitlements>"}

	}

	var groups []string
	groups, err = identityProvider.Groups(ctx)
	if err != nil {
		groups = []string{"<err-groups>"}

	}

	return gh.Nav(
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
    }`),
		gh.Class("relative z-10 w-auto"),
		gh.Div(gh.Class("relative"),
			gh.Ul(gh.Class("flex flex-1 p-1 space-x-1 list-none border rounded-md text-neutral-700 group border-neutral-200/80"),
				gh.Li(
					gh.Div(gp.Text("DataHub"), gh.Class(navBarTitleStyle)),
					gh.A(gh.Href("/collections"),
						gh.Class(navBarStyle),
						gp.Text("Collections")),
					gh.A(gh.Href("/labels"),
						gh.Class(navBarStyle),
						gp.Text("Labels")),
					gh.A(gh.Href("/sites"),
						gh.Class(navBarStyle),
						gp.Text("Sites")),
					gh.A(gh.Href("/profiles"),
						gh.Class(navBarStyle),
						gp.Text("Profiles")),
					gh.A(gh.Href("/docs"),
						gh.Rel("noopener noreferrer"),
						gh.Target("_blank"),
						gh.Class(navBarStyle),
						gp.Text("Docs")),
					gh.A(gh.Href("/api/v1"),
						gh.Rel("noopener noreferrer"),
						gh.Target("_blank"),
						gh.Class(navBarStyle),
						gp.Text("API")),
					gp.Raw(fmt.Sprintf(`<button
							:class="{ 'bg-neutral-100' : navigationMenu=='learn-more', 'hover:bg-neutral-100' : navigationMenu!='learn-more' }" @mouseover="navigationMenuOpen=true; navigationMenuReposition($el); navigationMenu='learn-more'" @mouseleave="navigationMenuLeave()" class="inline-flex items-center justify-center h-10 px-4 py-2 text-sm font-medium transition-colors rounded-md hover:text-neutral-900 focus:outline-none disabled:opacity-50 disabled:pointer-events-none bg-background hover:bg-neutral-100 group w-max">
							<span><div class="font-bold">%v</div></span>
						</button>`, username)),
				),
			)),
		gp.Raw(
			fmt.Sprintf(`
				<div x-ref="navigationDropdown" x-show="navigationMenuOpen"
					x-transition:enter="transition ease-out duration-100"
					x-transition:enter-start="opacity-0 scale-90"
					x-transition:enter-end="opacity-100 scale-100"
					x-transition:leave="transition ease-in duration-100"
					x-transition:leave-start="opacity-100 scale-100"
					x-transition:leave-end="opacity-0 scale-90"
					@mouseover="navigationMenuClearCloseTimeout()" @mouseleave="navigationMenuLeave()"
					class="absolute top-0 pt-3 duration-200 ease-out -translate-x-1/2 translate-y-11" x-cloak>

					<div class="flex justify-center w-auto h-auto overflow-hidden bg-white border rounded-md shadow-sm border-neutral-200/70">

						<div x-show="navigationMenu == 'learn-more'" class="flex items-stretch justify-center w-full p-6">
							<div class="w-72">
								<a href="#_" @click="navigationMenuClose()" class="block px-3.5 py-3 text-sm rounded hover:bg-neutral-100">
									<span class="block mb-1 font-medium text-black">Username</span>
									<span class="block font-light leading-5 opacity-50">%v</span>
								</a>
								<a href="#_" @click="navigationMenuClose()" class="block px-3.5 py-3 text-sm rounded hover:bg-neutral-100">
									<span class="block mb-1 font-medium text-black">Email</span>
									<span class="block font-light leading-5 opacity-50">%v</span>
								</a>
								<a href="#_" @click="navigationMenuClose()" class="block px-3.5 py-3 text-sm rounded hover:bg-neutral-100">
									<span class="block mb-1 font-medium text-black">Entitlements</span>
									<span class="block font-light leading-5 opacity-50">%v</span>
								</a>
								<a href="#_" @click="navigationMenuClose()" class="block px-3.5 py-3 text-sm rounded hover:bg-neutral-100">
									<span class="block mb-1 font-medium text-black">Groups</span>
									<span class="block font-light leading-5 opacity-50">%v</span>
								</a>
								<a href="%v" @click="navigationMenuClose()" class="block px-3.5 py-3 text-sm rounded hover:bg-neutral-100">
									<span class="block mb-1 font-medium text-black">Sign out</span>
								</a>
							</div>
						</div>

					</div>
				</div>`, username, email, entitlements, groups, signOutURL)),
	)
}
