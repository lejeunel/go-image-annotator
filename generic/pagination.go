package generic

import (
	"errors"
	"fmt"
	pag "github.com/vcraescu/go-paginator/v2"
)

func NewPaginator(adapter pag.Adapter, pageSize, maxPerPage, page int) (pag.Paginator, error) {

	if pageSize > maxPerPage {
		return nil, errors.New(fmt.Sprintf("Provided page size %d must be <= to %d", pageSize, maxPerPage))
	}

	if page < 1 {
		return nil, errors.New(fmt.Sprintf("Provided page %d must be >= to 1", page))
	}

	p := pag.New(adapter, maxPerPage)
	p.SetPage(page)
	return p, nil
}

type PaginationParams struct {
	Page     int `query:"page" minimum:"1" default:"1"`
	PageSize int `query:"pagesize"`
}

type PaginationMeta struct {
	Next         int   `json:"next,omitempty"`
	Previous     int   `json:"prev,omitempty"`
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

func NewPaginationMeta(p pag.Paginator) PaginationMeta {

	pagination := PaginationMeta{}

	hasNext, _ := p.HasNext()

	if hasNext {
		pagination.Next, _ = p.NextPage()
	}

	hasPrev, _ := p.HasPrev()

	if hasPrev {
		pagination.Previous, _ = p.PrevPage()
	}

	pagination.CurrentPage, _ = p.Page()
	pagination.TotalPages, _ = p.PageNums()
	pagination.TotalRecords, _ = p.Nums()

	return pagination

}
