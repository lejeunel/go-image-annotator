package clone

type Request struct {
	Source           string
	Destination      string
	DestinationGroup *string
	Deep             bool
}
