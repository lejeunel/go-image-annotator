package update

// Update label repo
type Repo interface {
	Update(Model) error
	Exists(string) (bool, error)
}
