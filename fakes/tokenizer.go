package fake

import (
	tk "github.com/lejeunel/go-image-annotator/entities/token"
)

type Tokenizer struct {
	GotToken    string
	ReturnHash  []byte
	ReturnValue string
}

func (t *Tokenizer) Hash(token string) []byte {
	t.GotToken = token
	return t.ReturnHash
}

func (t Tokenizer) Generate() (*tk.Token, error) {
	return &tk.Token{Value: t.ReturnValue, Hash: t.ReturnHash}, nil
}
