package server

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/server/routes"
	"net/http"
)

func ApiRequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := u.IdentityFromContext(r.Context())
		if user == nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func WebRequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := u.IdentityFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, rt.Login, http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
