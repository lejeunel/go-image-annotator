package token_generator

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
)

func AppendUserToToken(user string, token string) string {
	return user + ":" + token
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

type PersonalAccessToken struct {
	UserId   string
	APIToken string
}

func DecodeAndSplitPersonalAccessToken(input string) (*PersonalAccessToken, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("decoding token from base64: %w", err)
	}
	userId, apiToken, ok := strings.Cut(string(decoded), ":")
	if !ok {
		return nil, fmt.Errorf("splitting token")
	}
	return &PersonalAccessToken{userId, apiToken}, nil
}

type Pair struct {
	Value string
	Hash  []byte
}

type AuthGenerator interface {
	Generate() (*Pair, error)
	Hash(token string) []byte
	Verify(string, []byte) bool
}

type MyAuthGenerator struct {
	Length int
}

func New(length int) MyAuthGenerator {
	return MyAuthGenerator{Length: length}
}

func (g MyAuthGenerator) Generate() (*Pair, error) {
	buf := make([]byte, g.Length)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	token := base64.StdEncoding.EncodeToString(buf)
	sum := g.Hash(token)

	return &Pair{
		Value: token,
		Hash:  sum,
	}, nil
}
func (g MyAuthGenerator) Verify(token string, storedHash []byte) bool {
	computed := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare(computed[:], storedHash) == 1
}
func (g MyAuthGenerator) Hash(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}
