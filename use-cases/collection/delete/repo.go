package delete

type CollectionRepo interface {
	Delete(string) error
	Exists(string) (bool, error)
	IsPopulated(string) (*bool, error)
}

type GroupRepo interface {
	GroupOfCollection(string) (*string, error)
}
