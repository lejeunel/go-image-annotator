package clone

type Request struct {
	Source           string
	Destination      string
	DestinationGroup *string
	Deep             *bool
}

type Response struct {
	Id     string
	Issuer string
	Type   string
}
