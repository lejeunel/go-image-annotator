package builders

import (
	"bytes"
	"fmt"
	_ "strings"
	"text/template"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type NavBarActivatedItems struct {
	Home        bool
	Collections bool
	Labels      bool
	API         bool
}

var UserIcon = `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" aria-hidden="true" fill="currentColor"  class="w-1/2 h-1/2">
	<path fill-rule="evenodd" d="M7.5 6a4.5 4.5 0 119 0 4.5 4.5 0 01-9 0zM3.751 20.105a8.25 8.25 0 0116.498 0 .75.75 0 01-.437.695A18.683 18.683 0 0112 22.5c-2.786 0-5.433-.608-7.812-1.7a.75.75 0 01-.437-.695z" clip-rule="evenodd"/>
</svg>
`

type UserMenu struct {
	Icon     string
	UserName string
	Entries  []UserMenuEntry
}
type UserMenuEntry struct {
	Name string
	URL  string
}

func MakeLoginOrUserBadge(user *u.User) Node {
	if user == nil {
		return A(Href("/login/google"),
			Class("rounded-radius bg-primary border border-primary px-4 py-2 text-center text-sm font-medium tracking-wide text-on-primary hover:opacity-75 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary active:opacity-100 active:outline-offset-0 dark:bg-primary-dark dark:border-primary-dark dark:text-on-primary-dark dark:focus-visible:outline-primary-dark"),
			Text("Login"))
	}
	tUser := template.New("")
	template.Must(tUser.ParseFS(templatesFiles, "templates/user_badge.html"))
	var iconBuf bytes.Buffer
	Span(Class("flex size-10 items-center justify-center overflow-hidden rounded-full border border-outline bg-surface-alt text-on-surface/80 dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark/80"),
		Raw(UserIcon)).Render(&iconBuf)
	var buf bytes.Buffer
	entries := UserMenu{UserName: user.Id,
		Entries: []UserMenuEntry{{"Dashboard", "/user-dashboard"}, {"Sign Out", "/logout"}},
		Icon:    iconBuf.String()}
	tUser.ExecuteTemplate(&buf, "user_badge", entries)
	return Raw(buf.String())
}

func MakeRepoButton(url string) Node {
	icon := `
<svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="size-5" viewBox="0 0 16 16">
    <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8"></path>
</svg>
`
	return A(Href(url), Span(
		Class("text-onSurface dark:text-onSurfaceDark"),
		Target("_blank"),
		Raw(icon),
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
	moon := `<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z"></path>
			</svg>`
	sun := `<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"></path>
			</svg>`
	return Button(
		Attr("@click", "toggleDark()"),
		Attr("type", "button"),
		Class(`
			whitespace-nowrap hover:bg-gray-100 dark:hover:bg-gray-800 rounded-radius px-2 py-2 text-sm font-medium tracking-wide text-surface-dark
			transition hover:opacity-75 text-center focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-surface-dark
			active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed
			dark:text-surface dark:focus-visible:outline-surface
		`),
		Span(
			Attr("x-html", "darkMode ? `"+sun+"` : `"+moon+"`"),
			Attr(":class", "darkMode ? 'text-gray-300' : 'text-gray-700'"),
		),
	)
}
func MakeNavBar(isActivated ActivePage, repoURL string, apiPrefix string, user *u.User) Node {
	return Nav(
		Attr("x-on:click.away", "mobileMenuIsOpen = false"),
		Class("fixed top-0 z-30 hidden h-16 w-screen items-center justify-between border-outline px-10 py-2 backdrop-blur-xl md:flex dark:border-outline-dark bg-surface-alt/75 dark:bg-surface-dark-alt/75 border-b"),
		Aria("label", "penguin ui menu"),

		// --- Brand Logo ---
		A(
			Href("#"),
			Class("text-2xl font-bold text-on-surface-strong dark:text-on-surface-dark-strong"),
			Span(
				Text("Image"),
				Span(
					Class("text-primary dark:text-primary-dark"),
					Text("Annotator"),
				),
			),
		),

		// --- Desktop Menu ---
		Ul(
			Class("hidden items-center gap-4 md:flex"),
			Li(
				MakeMenuItem("Home", "/", isActivated == HomePageActive),
			),
			Li(
				MakeMenuItem("Collections", "/collections", isActivated == CollectionsPageActive),
			),
			Li(
				MakeMenuItem("Labels", "/labels", isActivated == LabelsPageActive),
			),
			Li(
				MakeMenuItem("API Docs", fmt.Sprintf("/%v/docs", apiPrefix), isActivated == APIDocsPageActive),
			),
			Li(MakeRepoButton(repoURL)),
			Li(
				DarkModeToggle(),
			),
			Li(MakeLoginOrUserBadge(user)),
		),
	)
}
