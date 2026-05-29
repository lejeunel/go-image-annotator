package label

type Label struct {
	Id          LabelId
	Name        string
	Description string
}

func NewLabel(id LabelId, name string, opts ...Option) Label {
	l := &Label{Id: id, Name: name}
	for _, opt := range opts {
		opt(l)
	}
	return *l
}

type Option func(*Label)

func WithDescription(d string) Option {
	return func(l *Label) {
		l.Description = d
	}
}
