package session

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	readusr "github.com/lejeunel/go-image-annotator/use-cases/user/read"
)

var UserIdKey = "user-id"

type SessionManager interface {
	Login(context.Context, string) error
	Logout(context.Context) error
	MiddleWare(next http.Handler) http.Handler
}

type TokenVerifier interface {
	Verify(token string, storedHash []byte) bool
}

type MySessionManager struct {
	*scs.SessionManager
	readusr.Repo
	TokenVerifier
}

func (m MySessionManager) MiddleWare(next http.Handler) http.Handler {
	return m.LoadAndSave(m.middlewareDecodeUserFromSessionId(m.LookForAPIToken(next)))
}

func (m MySessionManager) LookForAPIToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		bearerToken, ok := strings.CutPrefix(authHeader, "Bearer ")
		if ok && bearerToken != "" {
			token, err := tok.DecodeAndSplitToken(bearerToken)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := m.Repo.Find(token.UserId)
			if errors.Is(err, e.ErrNotFound) {
			}
			if match := m.TokenVerifier.Verify(token.APIToken, []byte(user.HashPAT)); !match {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), u.UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m MySessionManager) middlewareDecodeUserFromSessionId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := m.GetString(r.Context(), UserIdKey)
		if id == "" {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.Repo.Find(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), u.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (m MySessionManager) Logout(ctx context.Context) error {
	if err := m.SessionManager.RenewToken(ctx); err != nil {
		return err
	}
	return m.SessionManager.Clear(ctx)
}

func (m MySessionManager) Login(ctx context.Context, id string) error {
	errCtx := fmt.Errorf("logging in user %v", id)
	if _, err := m.Repo.Find(id); err != nil {
		return fmt.Errorf("%w: checking if user is registered: %w", errCtx, err)
	}

	if err := m.SessionManager.RenewToken(ctx); err != nil {
		return fmt.Errorf("%w: renewing token: %w", errCtx, err)
	}
	m.SessionManager.Put(ctx, UserIdKey, id)
	return nil
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf  bytes.Buffer
	code int
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	bw.code = code
}

func NewSQLiteSessionManager(db *sql.DB, repo readusr.Repo,
	verifier TokenVerifier) MySessionManager {
	store := sqlite3store.New(db)
	m := MySessionManager{SessionManager: scs.New(), Repo: repo,
		TokenVerifier: verifier}
	m.Store = store
	return m
}
