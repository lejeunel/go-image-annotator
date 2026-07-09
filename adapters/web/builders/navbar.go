package builders

import (
	"bytes"
	_ "strings"
	"text/template"

	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

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
	template.Must(tUser.ParseFS(templatesFiles, "templates/user_badge.html"))
	var iconBuf bytes.Buffer
	Raw(ic.UserIcon).Render(&iconBuf)
	var buf bytes.Buffer
	entries := UserMenu{UserName: user.Id,
		Entries: []UserMenuEntry{{"Dashboard", rt.UserDashboard}, {"Sign Out", rt.Logout}},
		Icon:    iconBuf.String()}
	tUser.ExecuteTemplate(&buf, "user_badge", entries)
	return Raw(buf.String())
}
func MakeRepoButton(url string) Node {
	return A(Href(url), Span(
		Class("text-onSurface dark:text-onSurfaceDark"),
		Target("_blank"),
		Raw(ic.GitHubIcon),
		Attr(":class", "darkMode ? 'text-gray-300' : 'text-gray-700'"),
	))
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
			whitespace-nowrap hover:bg-gray-100 dark:hover:bg-gray-800 rounded-radius px-2 py-2 text-sm font-medium tracking-wide text-surface-dark
			transition hover:opacity-75 text-center focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-surface-dark
			active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed
			dark:text-surface dark:focus-visible:outline-surface cursor-pointer
		`),
		Span(
			Attr("x-html", "darkMode ? `"+ic.SunIcon+"` : `"+ic.MoonIcon+"`"),
			Attr(":class", "darkMode ? 'text-gray-300' : 'text-gray-700'"),
		),
	)
}
func MakeNavBar(isActivated ActivePage, repoURL string, docsURL string, apiPrefix string, user u.User) Node {
	return Nav(
		Attr("x-on:click.away", "mobileMenuIsOpen = false"),
		Class("fixed top-0 z-30 hidden h-16 w-screen items-center justify-between border-outline px-10 py-2 backdrop-blur-xl md:flex dark:border-outline-dark bg-surface-alt/75 dark:bg-surface-dark-alt/75 border-b"),
		Aria("label", "penguin ui menu"),

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
				MakeMenuItem("Documentation", docsURL, false),
			),
			Li(
				MakeMenuItem("API", rt.APIDocs, isActivated == APIDocsPageActive),
			),
			Li(MakeRepoButton(repoURL)),
			Li(
				DarkModeToggle(),
			),
			Li(MakeUserBadge(user)),
		),
	)
}
