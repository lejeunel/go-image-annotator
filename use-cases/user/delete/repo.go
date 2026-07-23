package delete

type Repo interface {
	Exists(string) (bool, error)
	Delete(string) error
	CountAdmins() (int64, error)
}
