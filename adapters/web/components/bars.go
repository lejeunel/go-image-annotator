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

func MakeUserBadge(user u.User) Node {
	tUser := template.New("")
	template.Must(tUser.Parse(userBadgeTemplate))
	var iconBuf bytes.Buffer
	Raw(ic.UserCircle).Render(&iconBuf)
	var buf bytes.Buffer
	menu := UserMenu{UserName: user.Id, Icon: iconBuf.String()}
	menu.Entries = append(menu.Entries, UserMenuEntry{"Dashboard", rt.UserDashboard})
	if user.IsAdmin {
		menu.Entries = append(menu.Entries, UserMenuEntry{"Admin", rt.Admin})
	}
	menu.Entries = append(menu.Entries, UserMenuEntry{"Sign out", rt.Logout})
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
func MakeNavBar(isActivated ActivePage, repoURL string, docsURL string, apiPrefix string, user u.User) Node {
	return Nav(
		Attr("x-on:click.away", "mobileMenuIsOpen = false"),
		Class(
			"fixed top-0 z-30 hidden h-14 w-screen items-center justify-between border-outline px-10 py-2 backdrop-blur-xl md:flex dark:border-outline-dark bg-surface-alt/75 dark:bg-surface-dark-alt/75 border-b"),
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
				MakeMenuItem("Home", rt.Home, isActivated == HomePageActive),
			),
			Li(
				MakeMenuItem("Collections", rt.Collections, isActivated == CollectionsPageActive),
			),
			Li(
				MakeMenuItem("Labels", rt.Labels, isActivated == LabelsPageActive),
			),
			Li(
				A(
					Href(rt.APIDocs),
					Class("hover:text-primary"),
					Text("API"),
				),
			),
			Li(
				DarkModeToggle(),
			),
			Li(MakeUserBadge(user)),
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
		Class("flex fixed bottom-0 z-30 h-8 text-xs w-screen items-center justify-end border-t border-outline bg-surface-alt/75 px-10 backdrop-blur-xl dark:border-outline-dark dark:bg-surface-dark-alt/75"),
		Div(
			Class("flex items-center gap-2 text-gray-400 dark:text-gray-400 hover:text-gray-500 hover:dark:text-gray-500"),
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
