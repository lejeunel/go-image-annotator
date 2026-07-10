package components

type ActivePage int

const (
	CollectionsPageActive ActivePage = iota
	LabelsPageActive
	HomePageActive
	APIDocsPageActive
	NoPageActive
)
