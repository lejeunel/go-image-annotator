package builders

import "net/url"

type DynamicRowMode int

const (
	ModeView DynamicRowMode = iota
	ModeEdit
	ModeConfirmDelete
)

var modeName = map[DynamicRowMode]string{
	ModeView:          "view",
	ModeEdit:          "edit",
	ModeConfirmDelete: "confirm-delete",
}

func (m DynamicRowMode) String() string {
	return modeName[m]
}

type RowURL struct {
	Url       url.URL
	idArgName string
}

func NewRowURL(baseURL, idArgName string) RowURL {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	return RowURL{Url: *u, idArgName: idArgName}
}

func (b *RowURL) SetId(id string) *RowURL {
	q := b.Url.Query()
	q.Set(b.idArgName, id)
	b.Url.RawQuery = q.Encode()
	return b
}

func (b *RowURL) SetMode(m DynamicRowMode) *RowURL {
	q := b.Url.Query()
	q.Set("mode", m.String())
	b.Url.RawQuery = q.Encode()
	return b
}
func (b *RowURL) String() string {
	return b.Url.String()
}
