package components

type ActivePage int

const (
	CollectionsPageActive ActivePage = iota
	LabelsPageActive
	ProfilesPageActive
	HomePageActive
	APIDocsPageActive
	NoPageActive
)
