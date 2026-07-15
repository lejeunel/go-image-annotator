package ingester

type Hasher interface {
	Hash([]byte) string
}
