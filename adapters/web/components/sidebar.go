package components

import (
	"bytes"
	"io"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type SidebarEntry struct {
	icon     string
	label    string
	url      string
	isActive bool
}

func NewSidebarEntry(icon, label, url string) SidebarEntry {
	return SidebarEntry{icon, label, url, false}
}
func (e *SidebarEntry) SetActive(value bool) {
	e.isActive = value
}

func (e SidebarEntry) Render(w io.Writer) {
	aClass := "flex items-center rounded-radius gap-2 px-2 py-1.5 text-sm font-medium text-on-surface underline-offset-2 hover:bg-primary/5 hover:text-on-surface-strong focus-visible:underline focus:outline-hidden dark:text-on-surface-dark dark:hover:bg-primary-dark/5 dark:hover:text-on-surface-dark-strong"
	if e.isActive {
		aClass = "flex items-center rounded-radius gap-2 bg-primary/10 px-2 py-1.5 text-sm font-medium text-on-surface-strong underline-offset-2 focus-visible:underline focus:outline-hidden dark:bg-primary-dark/10 dark:text-on-surface-dark-strong"
	}
	A(Href(e.url), Class(aClass), Raw(e.icon), Span(Text(e.label))).Render(w)
}

type Sidebar struct {
	Title        string
	Entries      map[string]SidebarEntry
	entriesNames []string
}

func NewSidebar(title string) Sidebar {
	return Sidebar{Title: title, Entries: make(map[string]SidebarEntry)}
}
func (s *Sidebar) Activate(name string) {
	entry, ok := s.Entries[name]
	if !ok {
		return
	}
	entry.isActive = true
	s.Entries[name] = entry
}

func (s *Sidebar) AddEntry(name, icon, url string) *Sidebar {
	s.Entries[name] = NewSidebarEntry(icon, name, url)
	s.entriesNames = append(s.entriesNames, name)
	return s
}
func (s Sidebar) Render(w io.Writer) {
	var buf bytes.Buffer
	Span(Class("ml-2 w-fit text-xl font-bold text-on-surface-strong dark:text-on-surface-dark-strong"),
		Text(s.Title)).Render(&buf)

	for _, n := range s.entriesNames {
		s.Entries[n].Render(&buf)
	}
	nodes := Div(Class("flex flex-col gap-2 overflow-y-auto pb-6"), Raw(buf.String()))
	nodes.Render(w)
}
