package fake

type Hasher struct {
	Sum_ []byte
}

func (f *Hasher) Write(p []byte) (int, error) {
	return len(p), nil
}

func (f *Hasher) Sum(b []byte) []byte {
	return append(b, f.Sum_...)
}

func (f *Hasher) Reset() {}

func (f *Hasher) Size() int {
	return len(f.Sum_)
}

func (f *Hasher) BlockSize() int {
	return 1
}
