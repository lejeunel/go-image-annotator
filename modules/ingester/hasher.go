package ingest

type Hasher interface {
	Hash([]byte) string
}
