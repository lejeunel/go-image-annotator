package session

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	readusr "github.com/lejeunel/go-image-annotator/use-cases/user/read"
	"net/http"
)

var UserIdKey = "user-id"

type SessionManager interface {
	Login(context.Context, string) error
	Logout(context.Context) error
	MiddleWare(next http.Handler) http.Handler
}

type MySessionManager struct {
	*scs.SessionManager
	readusr.Repo
}

func (m MySessionManager) MiddleWare(next http.Handler) http.Handler {
	return m.LoadAndSave(m.LoadAndSaveWithHeader(m.appendUserIdentityToContext(next)))
}

func (m MySessionManager) LoadAndSaveWithHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerKey := "X-Session"
		headerKeyExpiry := "X-Session-Expiry"

		ctx, err := m.Load(r.Context(), r.Header.Get(headerKey))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		bw := &bufferedResponseWriter{ResponseWriter: w}
		sr := r.WithContext(ctx)
		next.ServeHTTP(bw, sr)

		if m.Status(ctx) == scs.Modified {
			token, expiry, err := m.Commit(ctx)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.Header().Set(headerKey, token)
			w.Header().Set(headerKeyExpiry, expiry.Format(http.TimeFormat))
		}

		if bw.code != 0 {
			w.WriteHeader(bw.code)
		}
		w.Write(bw.buf.Bytes())
	})
}

func (m MySessionManager) appendUserIdentityToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := m.Repo.Find(m.GetString(r.Context(), UserIdKey))
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
	if err := m.SessionManager.RenewToken(ctx); err != nil {
		return err
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

func NewSQLiteSessionManager(db *sql.DB, repo readusr.Repo) MySessionManager {
	store := sqlite3store.New(db)
	m := MySessionManager{SessionManager: scs.New(), Repo: repo}
	m.Store = store
	return m
}
