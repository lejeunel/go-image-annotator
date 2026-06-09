package update

type GroupRepo interface {
	GroupOfCollection(string) (*string, error)
}

type CollectionRepo interface {
	Update(Model) error
	Exists(string) (bool, error)
}
