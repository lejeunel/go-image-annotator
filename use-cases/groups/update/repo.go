package update

type Repo interface {
	Update(Model) error
	Exists(string) (bool, error)
	GroupOfCollection(string) (*string, error)
}
