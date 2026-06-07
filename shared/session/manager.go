package session

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	i "github.com/lejeunel/go-image-annotator/entities/identity"
	"net/http"
)

type SessionManager interface {
	Login(context.Context, string) error
	Logout(context.Context) error
	MiddleWare(next http.Handler) http.Handler
}

type MySessionManager struct {
	*scs.SessionManager
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

// Append a User record in the current request context.
// This looks in the session manager's store for a user-id.
// If it exists, it queries the user repository and append
// the user record to current context
// TODO This appends a hard-coded user for now.
func (m MySessionManager) appendUserIdentityToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := m.GetString(r.Context(), "user-id")
		if userId == "" {
			next.ServeHTTP(w, r)
			return
		}
		identity := &i.Identity{Id: userId, Groups: []string{"group-0", "group-1"},
			Roles: []string{"reader", "writer"}}
		ctx := context.WithValue(r.Context(), i.UserKey, identity)
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
	m.SessionManager.Put(ctx, "user-id", id)
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

func NewSQLiteSessionManager(db *sql.DB) MySessionManager {
	store := sqlite3store.New(db)
	m := MySessionManager{scs.New()}
	m.Store = store
	return m
}
