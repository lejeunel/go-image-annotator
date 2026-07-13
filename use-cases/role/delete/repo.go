package delete

type Repo interface {
	Delete(string) error
	Exists(string) (*bool, error)
	IsAssigned(string) (*bool, error)
}
