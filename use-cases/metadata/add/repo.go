package add

type CollectionRepo interface {
	GetGroup(string) (*string, error)
}

type MetaDataRepo interface {
	Add(string, any) error
	KeyExists(string) (*bool, error)
}
