package update

type Repo interface {
	Update(Model) error
	Exists(string) (bool, error)
	Group(string) (*string, error)
}
