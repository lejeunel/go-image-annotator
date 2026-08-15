package components

import (
	"bytes"
	_ "strings"
	"text/template"

	_ "embed"

	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	g "github.com/lejeunel/go-image-annotator/globals"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed templates/user_badge.html
var userBadgeTemplate string

type NavBarActivatedItems struct {
	Home        bool
	Collections bool
	Labels      bool
	API         bool
}

type UserMenu struct {
	Icon     string
	UserName string
	Entries  []UserMenuEntry
}
type UserMenuEntry struct {
	Name string
	URL  string
}

func MakeUserBadge(user u.User, dashboardURL string) Node {
	tUser := template.New("")
	template.Must(tUser.Parse(userBadgeTemplate))
	var iconBuf bytes.Buffer
	Raw(ic.UserCircle).Render(&iconBuf)
	var buf bytes.Buffer
	menu := UserMenu{UserName: user.Id, Icon: iconBuf.String()}
	menu.Entries = append(menu.Entries, UserMenuEntry{"Dashboard", dashboardURL})
	if user.IsAdmin() {
		menu.Entries = append(menu.Entries, UserMenuEntry{"Admin", rt.AdminUrl})
	}
	menu.Entries = append(menu.Entries, UserMenuEntry{"Sign out", rt.LogoutUrl})
	tUser.ExecuteTemplate(&buf, "user_badge", menu)
	return Raw(buf.String())
}

func MakeRepoButton(repoName string, currentVersion, url string) Node {
	return A(
		Target("_blank"),
		Href(url),
		Div(
			Class("flex items-center gap-1"),
			Span(Raw(ic.GitHub)),
			Span(Text(repoName)),
			Span(Text(currentVersion)),
		),
	)
}

func MakeMenuItem(name string, url string, activated bool) Node {
	class := "font-medium text-on-surface underline-offset-2 hover:text-primary focus:outline-hidden focus:underline dark:text-on-surface-dark dark:hover:text-primary-dark"
	if activated {
		class = "font-bold text-primary underline-offset-2 hover:text-primary focus:outline-hidden focus:underline dark:text-primary-dark dark:hover:text-primary-dark"
	}

	return A(
		Href(url),
		Aria("current", "page"),
		Span(Class(class), Text(name)),
	)
}

func DarkModeToggle() Node {
	return Button(
		Attr("@click", "toggleDark()"),
		Attr("type", "button"),
		Class(`
			whitespace-nowrap hover:bg-gray-100 dark:hover:bg-gray-800 rounded-radius px-1 py-2 text-sm font-medium tracking-wide text-surface-dark
			transition hover:opacity-75 text-center focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-surface-dark
			active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed
			dark:text-surface dark:focus-visible:outline-surface cursor-pointer
		`),
		Span(
			Attr("x-html", "darkMode ? `"+ic.Sun+"` : `"+ic.Moon+"`"),
			Attr(":class", "darkMode ? 'text-gray-300' : 'text-gray-700'"),
		),
	)
}

func MakeSearchButton() Node {

	return Div(Class("bg-surface-alt p-0.5 dark:bg-surface-dark-alt"),
		Button(
			Class("btn-search flex h-8 w-full cursor-pointer items-center justify-between border-outline bg-surface px-2 font-light transition-all duration-200 dark:border-outlineDark dark:bg-surface-dark rounded-lg border"),
			Attr(`x-on:click="showSearch=true, $dispatch('searchModalOpened')"`),
			Attr(`x-bind:class="['rounded-lg', 'border']"`),
			Div(Class("flex items-center gap-2"),
				Raw(`
	        <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-5 w-5 text-on-surface dark:text-on-surface-dark">
	            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"></path>
	        </svg>

			`),
				// Text("Search"),
			),
			Div(
				Attr(`x-data="{ os: detectOS() }"`),
				Class("flex items-center gap-1 text-xs text-on-surface-strong dark:text-on-surface-dark"),
				Raw(`
	        <svg x-show="os === 'Mac OS'" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="currentColor" viewBox="0 0 16 16" style="display: none;">
	            <path d="M3.5 2A1.5 1.5 0 0 1 5 3.5V5H3.5a1.5 1.5 0 1 1 0-3zM6 5V3.5A2.5 2.5 0 1 0 3.5 6H5v4H3.5A2.5 2.5 0 1 0 6 12.5V11h4v1.5a2.5 2.5 0 1 0 2.5-2.5H11V6h1.5A2.5 2.5 0 1 0 10 3.5V5H6zm4 1v4H6V6h4zm1-1V3.5A1.5 1.5 0 1 1 12.5 5H11zm0 6h1.5a1.5 1.5 0 1 1-1.5 1.5V11zm-6 0v1.5A1.5 1.5 0 1 1 3.5 11H5z"></path>
	        </svg>
					`),
				Span(Attr(`x-show="os === 'Windows' || os === 'Linux'" aria-hidden="true"`), Div(Class("pl-1"), Text("Ctrl +"))),
				Span(Attr(`x-show="os === 'Mac OS' || os === 'Linux' || os === 'Windows'" aria-hidden="true"`), Text("K")),
			),
		),
	)
	// <button type="button" aria-label="Search" class="btn-search flex h-10 w-full cursor-pointer items-center justify-between border-outline bg-surfaceAlt p-2 px-4 font-light transition-all duration-200 dark:border-outlineDark dark:bg-surfaceDarkAlt rounded-lg border"
	// x-on:click="showSearch=true, $dispatch('searchModalOpened')" x-bind:class="[
	//         'rounded' + $store.border.radius,
	//         $store.border.border ? 'border' : 'border-none'
	//     ]">
	//     <div class="flex items-center gap-2">
	//         <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-5 w-5 text-onSurface dark:text-onSurfaceDark">
	//             <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"></path>
	//         </svg>
	//         Search
	//     </div>
	//     <div x-data="{ os: detectOS() }" class="flex items-center gap-1 text-sm text-onSurfaceStrong dark:text-onSurfaceDark">
	//         <svg x-show="os === 'Mac OS'" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="currentColor" viewBox="0 0 16 16" style="display: none;">
	//             <path d="M3.5 2A1.5 1.5 0 0 1 5 3.5V5H3.5a1.5 1.5 0 1 1 0-3zM6 5V3.5A2.5 2.5 0 1 0 3.5 6H5v4H3.5A2.5 2.5 0 1 0 6 12.5V11h4v1.5a2.5 2.5 0 1 0 2.5-2.5H11V6h1.5A2.5 2.5 0 1 0 10 3.5V5H6zm4 1v4H6V6h4zm1-1V3.5A1.5 1.5 0 1 1 12.5 5H11zm0 6h1.5a1.5 1.5 0 1 1-1.5 1.5V11zm-6 0v1.5A1.5 1.5 0 1 1 3.5 11H5z"></path>
	//         </svg>
	//         <span x-show="os === 'Windows' || os === 'Linux'" aria-hidden="true">
	//             Ctrl +
	//         </span>
	//         <span x-show="os === 'Mac OS' || os === 'Linux' || os === 'Windows'" aria-hidden="true">
	//             K
	//         </span>
	//     </div>
	// </button>

}

func MakeNavBar(
	isActivated ActivePage,
	repoURL string,
	docsURL string,
	apiPrefix string,
	user u.User,
	userDashboardURL string,
) Node {
	return Nav(
		Attr("x-on:click.away", "mobileMenuIsOpen = false"),
		Class(
			"fixed top-0 z-30 hidden h-14 w-screen items-center justify-between border-outline px-10 py-2 backdrop-blur-xl md:flex dark:border-outline-dark bg-surface-alt/75 dark:bg-surface-dark-alt/75 border-b",
		),
		Aria("label", "ui menu"),

		A(
			Href("/"),
			Class("text-2xl font-bold text-on-surface-strong dark:text-on-surface-dark-strong"),
			Span(
				Text("Image"),
				Span(
					Class("text-primary dark:text-primary-dark"),
					Text("Annotator"),
				),
			),
		),

		Ul(
			Class("hidden items-center gap-4 md:flex"),
			Li(
				MakeSearchButton(),
			),
			Li(
				MakeMenuItem("Home", rt.HomePageUrl, isActivated == HomePageActive),
			),
			Li(
				MakeMenuItem(
					"Collections",
					rt.CollectionsUrl,
					isActivated == CollectionsPageActive,
				),
			),
			Li(
				MakeMenuItem("Labels", rt.LabelsUrl, isActivated == LabelsPageActive),
			),
			Li(
				A(
					Href(rt.APIDocsUrl),
					Class("hover:text-primary"),
					Text("API"),
				),
			),
			Li(
				DarkModeToggle(),
			),
			Li(MakeUserBadge(user, userDashboardURL)),
		),
	)
}

func MakeDocsButton() Node {
	return A(
		Target("_blank"),
		Href(g.DocsURL),
		Div(
			Class("flex items-center gap-1"),
			Span(Raw(ic.Book)),
			Span(Text("Docs")),
		),
	)
}

func MakeFooter(currentVersion g.Info) Node {
	return Footer(
		Class(
			"flex fixed bottom-0 z-30 h-8 text-xs w-screen items-center justify-end border-t border-outline bg-surface-alt/75 px-5 py-3 backdrop-blur-xl dark:border-outline-dark dark:bg-surface-dark-alt/75",
		),
		Div(
			Class(
				"flex items-center gap-2 py-1 text-gray-400 dark:text-gray-400 hover:text-gray-500 hover:dark:text-gray-500",
			),
			Div(
				MakeDocsButton(),
			),
			Div(Class("h-5 w-px bg-gray-800 dark:bg-gray-400")),
			Div(
				MakeRepoButton(g.PackageName, currentVersion.Version, g.RepoURL),
			),
		),
	)
}
