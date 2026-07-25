package fake

import (
	tk "github.com/lejeunel/go-image-annotator/entities/token"
)

type Tokenizer struct {
	Hashed      string
	ReturnHash  []byte
	ReturnValue string
	FailVerify  bool
}

func (t *Tokenizer) Hash(token string) []byte {
	t.Hashed = token
	return t.ReturnHash
}

func (t Tokenizer) Generate() (*tk.Token, error) {
	return &tk.Token{Value: t.ReturnValue, Hash: t.ReturnHash}, nil
}

func (t *Tokenizer) Verify(token string, hash []byte) bool {
	t.Hashed = token
	if t.FailVerify {
		return false
	}
	return true
}
