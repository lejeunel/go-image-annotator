package token

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	tk "github.com/lejeunel/go-image-annotator/entities/token"
	"strings"
)

type TokenVerifier interface {
	Verify(token string, storedHash []byte) bool
}

type TokenHasher interface {
	Hash(token string) []byte
}

type TokenHashVerifier interface {
	Verify(token string, storedHash []byte) bool
	Hash(token string) []byte
}

func AppendUserToToken(user string, token string) string {
	return user + ":" + token
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

type APIToken struct {
	UserId   string
	APIToken string
}

func DecodeAndSplitPersonalAccessToken(input string) (*APIToken, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("decoding token from base64: %w", err)
	}
	userId, apiToken, ok := strings.Cut(string(decoded), ":")
	if !ok {
		return nil, fmt.Errorf("splitting token")
	}
	return &APIToken{userId, apiToken}, nil
}

type Interface interface {
	Generate() (*tk.Token, error)
	Verify(string, []byte) bool
	TokenHasher
}

type TokenService struct {
	Length int
	TokenHasher
}

func New(length int) TokenService {
	return TokenService{
		Length:      length,
		TokenHasher: NewSHA256Hasher()}
}

func (g TokenService) Generate() (*tk.Token, error) {
	buf := make([]byte, g.Length)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	token := base64.StdEncoding.EncodeToString(buf)
	sum := g.Hash(token)

	return &tk.Token{
		Value: token,
		Hash:  sum,
	}, nil
}
func (g TokenService) Verify(token string, storedHash []byte) bool {
	computed := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare(computed[:], storedHash) == 1
}

type SHA256Hasher struct{}

func NewSHA256Hasher() SHA256Hasher {
	return SHA256Hasher{}
}

func (h SHA256Hasher) Hash(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}
