package session

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	readusr "github.com/lejeunel/go-image-annotator/use-cases/user/find"
)

var UserIdKey = "user-id"

type SessionManager interface {
	FinishOAuthLogin(context.Context, string) error
	Logout(context.Context) error
	PasswordLogin(context.Context, string, string) error
}

type MySessionManager struct {
	*scs.SessionManager
	readusr.Repo
	tk.TokenVerifier
}

func (m MySessionManager) AuthCookiesMiddleWare(next http.Handler) http.Handler {
	return m.LoadAndSave(m.AuthFromSessionId(next))
}

func (m MySessionManager) AuthBearerMiddleWare(next http.Handler) http.Handler {
	return m.LoadAndSave(m.AuthBearerToken(next))
}
func (m MySessionManager) fetchUserFromBearerToken(bearerToken string) (*u.User, error) {
	errCtx := fmt.Errorf("inferring user's identity from bearer token")
	token, err := tk.DecodeAndSplitPersonalAccessToken(bearerToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCtx, err)
	}
	user, err := m.Repo.Find(token.UserId)
	if err != nil {
		return nil, fmt.Errorf("%w: fetching user: %w", errCtx, err)
	}
	if match := m.TokenVerifier.Verify(token.APIToken, []byte(user.HashPAT)); !match {
		return nil, fmt.Errorf("%w: verifying token: %w", errCtx, err)
	}
	return user, nil
}
func (m MySessionManager) AuthBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		bearerToken, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || bearerToken == "" {
			http.Error(w, "got invalid bearer token", http.StatusUnauthorized)
			return
		}

		user, err := m.fetchUserFromBearerToken(bearerToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := u.AppendUserToContext(r.Context(), *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (m MySessionManager) AuthFromSessionId(next http.Handler) http.Handler {
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
		ctx := u.AppendUserToContext(r.Context(), *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (m MySessionManager) Logout(ctx context.Context) error {
	if err := m.SessionManager.RenewToken(ctx); err != nil {
		return err
	}
	return m.SessionManager.Clear(ctx)
}
func (m MySessionManager) FinishOAuthLogin(ctx context.Context, id string) error {
	errCtx := fmt.Errorf("logging in user %v", id)
	if _, err := m.Repo.Find(id); err != nil {
		return fmt.Errorf("%w: checking if user is registered: %w", errCtx, err)
	}

	if err := m.initSession(ctx, id); err != nil {
		return err
	}
	return nil
}
func (m MySessionManager) PasswordLogin(ctx context.Context, email, password string) error {
	errCtx := fmt.Errorf("logging in user %v using password method", email)
	user, err := m.Repo.Find(email)
	if err != nil {
		return fmt.Errorf("%w: fetching user from email: %w", errCtx, err)
	}

	if !m.TokenVerifier.Verify(password, []byte(user.HashPassword)) {
		return fmt.Errorf("%w: matching password: %w", errCtx, e.ErrPasswordMismatch)
	}
	if err := m.initSession(ctx, user.Id); err != nil {
		return err
	}
	return nil
}

func (m MySessionManager) initSession(ctx context.Context, id string) error {
	if err := m.SessionManager.RenewToken(ctx); err != nil {
		return fmt.Errorf("initializing session: renewing token %w", err)
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
	verifier tk.TokenVerifier) MySessionManager {
	store := sqlite3store.New(db)
	m := MySessionManager{SessionManager: scs.New(), Repo: repo,
		TokenVerifier: verifier}
	m.Store = store
	return m
}
