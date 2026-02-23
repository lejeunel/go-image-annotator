package locations

type FilterArgs struct {
	Group      *string `json:"group" doc:"Group"`
	Collection *string `json:"collection" doc:"Collection name"`
}

func (f *FilterArgs) GetGroup() string {
	if f.Group != nil {
		return *f.Group
	}
	return ""
}

func (f *FilterArgs) GetCollection() string {
	if f.Collection != nil {
		return *f.Collection
	}
	return ""
}

func NewSiteByGroupFilter(group string) FilterArgs {
	return FilterArgs{Group: &group}
}

type OrderingArgs struct {
	Name bool
}

var SiteAlphabeticalOrdering = OrderingArgs{Name: true}

type FilterOption func(*FilterArgs)

func WithGroup(group string) FilterOption {
	return func(f *FilterArgs) {
		f.Group = &group
	}
}

func WithCollection(collection string) FilterOption {
	return func(f *FilterArgs) {
		f.Collection = &collection
	}
}

func NewSiteFilter(opts ...FilterOption) *FilterArgs {
	filter := &FilterArgs{}
	for _, opt := range opts {
		opt(filter)
	}
	return filter
}
