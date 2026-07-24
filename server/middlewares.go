package server

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

// Middleware is the standard net/http middleware signature.
type Middleware func(http.Handler) http.Handler

// Chain composes multiple middlewares into a single one.
// They are applied in the order given, i.e. Chain(A, B, C)(h) == A(B(C(h))),
// meaning A runs first on the way in.
func Chain(mws ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		h := final
		for i := len(mws) - 1; i >= 0; i-- {
			h = mws[i](h)
		}
		return h
	}
}

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
			http.Redirect(w, r, rt.LoginPageUrl, http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
