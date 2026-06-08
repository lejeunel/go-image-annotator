package delete

type Repo interface {
	Delete(string) error
	Exists(string) (bool, error)
	IsPopulated(string) (*bool, error)
	GroupOfCollection(string) (*string, error)
}
