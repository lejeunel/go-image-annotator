package components

import (
	_ "embed"
	"fmt"
	au "github.com/lejeunel/go-image-annotator/modules/token"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	"html/template"
	"io"
	"net/http"
)

//go:embed templates/api_token_display.html
var ApiTokenDisplay string

//go:embed templates/api_token_frame.html
var ApiTokenFrame string

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
	t := template.New("")
	template.Must(t.Parse(ApiTokenDisplay))
	t.Execute(p.Writer,
		TokenData{Token: au.Base64Encode(au.AppendUserToToken(resp.Id,
			resp.PersonalAccessToken))})
}
func NewAPITokenPresenter(w http.ResponseWriter) APITokenPresenter {
	return APITokenPresenter{w}
}
