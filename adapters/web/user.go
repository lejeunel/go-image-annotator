package web

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	tg "github.com/lejeunel/go-image-annotator/app/token-generator"
	p "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
)

func (s *Server) UserDashboard(w http.ResponseWriter, r *http.Request) {
	p := s.PageBuilder
	p.SetUserIdentityFromContext(r.Context())
	udb := s.UserDashboardBuilder.SetUserIdentityFromContext(r.Context())
	p.SetTitle("User Dashboard")
	p.SetContent(udb.Build())
	p.Render(w)

}

type APITokenPresenter struct {
	Writer io.Writer
}

func (p APITokenPresenter) Error(err error) {
	fmt.Println(err.Error())
}

type TokenData struct {
	Token string
}

func (p APITokenPresenter) Success(resp rt.Response) {
	t, err := template.ParseFS(templatesFiles, "templates/api_token_display.html")
	if err != nil {
		p.Writer.Write([]byte(err.Error()))
	}
	t.Execute(p.Writer,
		TokenData{Token: tg.Base64Encode(tg.AppendUserToToken(resp.Id,
			resp.PersonalAccessToken))})
}
func NewAPITokenPresenter(w http.ResponseWriter) APITokenPresenter {
	return APITokenPresenter{w}
}
func (s *Server) NewAPIToken(w http.ResponseWriter, r *http.Request) {
	user := p.IdentityFromContext(r.Context())
	if user == nil {
		http.Error(w, "failed getting user identity", http.StatusForbidden)
	}
	s.Interactors.User.RenewToken.Execute(r.Context(),
		rt.Request{Id: user.Id}, NewAPITokenPresenter(w))
}
