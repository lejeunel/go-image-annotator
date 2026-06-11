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

type IdentifiedToken struct {
	UserId   string
	APIToken string
}

func DecodeAndSplitToken(input string) (*IdentifiedToken, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("decoding token from base64: %w", err)
	}
	userId, apiToken, ok := strings.Cut(string(decoded), ":")
	if !ok {
		return nil, fmt.Errorf("splitting token")
	}
	return &IdentifiedToken{userId, apiToken}, nil
}

type TokenPair struct {
	Token string
	Hash  []byte
}

type TokenGenerator interface {
	Generate() (*TokenPair, error)
	Hash(token string) []byte
	Verify(string, []byte) bool
}

type MyTokenGenerator struct {
	Length int
}

func NewTokenGenerator(length int) MyTokenGenerator {
	return MyTokenGenerator{Length: length}
}

func (g MyTokenGenerator) Generate() (*TokenPair, error) {
	buf := make([]byte, g.Length)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(buf)
	sum := g.Hash(token)

	return &TokenPair{
		Token: token,
		Hash:  sum,
	}, nil
}
func (g MyTokenGenerator) Verify(token string, storedHash []byte) bool {
	computed := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare(computed[:], storedHash) == 1
}

func (g MyTokenGenerator) Hash(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}
