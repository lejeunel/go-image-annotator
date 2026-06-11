package builders

import (
	"fmt"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	"html/template"
	"io"
	"net/http"
)

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
		TokenData{Token: tok.Base64Encode(tok.AppendUserToToken(resp.Id,
			resp.PersonalAccessToken))})
}
func NewAPITokenPresenter(w http.ResponseWriter) APITokenPresenter {
	return APITokenPresenter{w}
}
