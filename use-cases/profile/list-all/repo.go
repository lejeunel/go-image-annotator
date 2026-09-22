package list

type Repo interface {
	ListAll() ([]string, error)
}
