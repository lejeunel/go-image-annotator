package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"net/http"
)

func ForgotPassword(builder b.ForgotPasswordBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		builder.Render(w)
	}
}
