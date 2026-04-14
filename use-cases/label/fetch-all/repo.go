package fetchall

type Repo interface {
	FetchAll() ([]string, error)
	Count() (int64, error)
}
