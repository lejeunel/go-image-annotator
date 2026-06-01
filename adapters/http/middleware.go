package http

import (
	"context"
	"github.com/lejeunel/go-image-annotator/entities/principal"
	"net/http"
)

func DummyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := principal.Principal{
			Id:     "anonymous",
			Email:  "anonymous@mail.com",
			Groups: []string{"the-group"},
			Roles:  []string{"admin"},
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "principal", p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
