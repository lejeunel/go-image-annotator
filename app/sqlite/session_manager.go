package sqlite

import (
	"database/sql"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
	sm "github.com/lejeunel/go-image-annotator/shared/session"
	readusr "github.com/lejeunel/go-image-annotator/use-cases/user/find"
)

func NewSessionManager(db *sql.DB, repo readusr.UserStore,
	verifier tk.TokenVerifier,
) sm.SessionManager {
	store := sqlite3store.New(db)
	m := sm.SessionManager{
		SessionManager: scs.New(), UserStore: repo,
		TokenVerifier: verifier,
	}
	m.Store = store
	return m
}
