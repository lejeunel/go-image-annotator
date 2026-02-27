package generic

import (
	e "datahub/errors"
	"fmt"
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	"math"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strconv"
)

type PaginationParams struct {
	Page     int64 `query:"page" minimum:"1" default:"1"`
	PageSize int   `query:"pagesize" default:"2"`
}

var OneItemPaginationParams = PaginationParams{Page: 1, PageSize: 1}

func NewPaginationParamsFromURL(url url.URL, maxPageSize int) (*PaginationParams, error) {
	pageint, err := strconv.Atoi(url.Query().Get("page"))
	if err != nil {
		pageint = 1
	}

	var pageSize int
	pageSize, err = strconv.Atoi(url.Query().Get("pagesize"))
	if err != nil {
		pageSize = maxPageSize
	}
	return NewPaginationParams(int64(pageint), pageSize, maxPageSize)

}

func NewPaginationParams(page int64, pageSize int, maxPageSize int) (*PaginationParams, error) {
	paginationParams := &PaginationParams{Page: page, PageSize: pageSize}
	if err := paginationParams.Validate(maxPageSize); err != nil {
		return nil, err
	}
	return paginationParams, nil
}

func (p *PaginationParams) Validate(maxPageSize int) error {
	if p.PageSize > maxPageSize {
		return fmt.Errorf("pagesize=%v exceeds maximum pagesize=%v: %w", p.PageSize, maxPageSize, e.ErrPagination)
	}

	if p.PageSize < 1 {
		return fmt.Errorf("pagesize=%v is inferior to 1: %w", p.PageSize, e.ErrPagination)
	}

	if p.Page < 1 {
		return fmt.Errorf("page=%v is inferior to 1: %w", p.Page, e.ErrPagination)
	}
	return nil

}

func (p PaginationParams) Limit() int {
	return p.PageSize
}

func (p PaginationParams) Offset() int64 {
	return int64(p.PageSize) * (p.Page - 1)
}

type PaginationMeta struct {
	Next         int64 `json:"next,omitempty"`
	Previous     int64 `json:"prev,omitempty"`
	CurrentPage  int64 `json:"current_page"`
	PageSize     int64 `json:"page_size"`
	TotalPages   int64 `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

func NewPaginationMeta(currentPage, totalRecords, pageSize int64) PaginationMeta {

	totalPages := int64(math.Ceil(float64(totalRecords) / float64(pageSize)))

	next := int64(0)
	if currentPage < totalPages {
		next = currentPage + 1
	}

	prev := int64(0)
	if currentPage > 1 {
		prev = currentPage - 1
	}

	return PaginationMeta{CurrentPage: currentPage,
		TotalPages: totalPages, TotalRecords: totalRecords,
		Next:     next,
		Previous: prev,
		PageSize: pageSize}

}

func GetPage(r http.Request) int64 {
	var page int
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	return int64(page)
}

var UrlClass = "font-medium text-blue-600 dark:text-blue-500 hover:underline"

func FixedSizePageArray(current, last, size int) []int {
	var pages []int
	left := max(1, current-1)
	right := min(current+1, last)
	nAdded := 0
	for {
		addedOne := false

		if left > 0 {
			pages = append(pages, left)
			left -= 1
			nAdded += 1
			addedOne = true
		}
		if nAdded >= size {
			break
		}
		if right < last+1 {
			pages = append(pages, right)
			right += 1
			nAdded += 1
			addedOne = true
		}
		if (nAdded >= size) || (addedOne == false) {
			break
		}
	}

	if !slices.Contains(pages, current) {
		pages = append(pages, current)
	}
	sort.Slice(pages, func(i, j int) bool {
		return pages[i] < pages[j]
	})

	return pages
}

func MakePaginationViewer(p PaginationMeta, baseURL url.URL, limitPages, paginationWidgetSize int) (*gp.Node, error) {

	begin := 1 + (int(p.CurrentPage)-1)*limitPages
	end := min(int(p.CurrentPage)*limitPages, int(p.TotalRecords))
	numEntries := int(p.TotalRecords)

	pages := FixedSizePageArray(int(p.CurrentPage), int(p.TotalPages), paginationWidgetSize)

	prevURL := baseURL
	nextURL := baseURL

	prevValues := prevURL.Query()
	prevValues.Set("page", strconv.Itoa(int(p.CurrentPage)-1))
	nextValues := nextURL.Query()
	nextValues.Set("page", strconv.Itoa(int(p.CurrentPage)+1))

	nextURL.RawQuery = nextValues.Encode()
	prevURL.RawQuery = prevValues.Encode()

	midURL := baseURL

	paginator := gh.Nav()
	if p.TotalPages > 1 {
		paginator = gh.Nav(
			gh.Ul(gh.Class("inline-flex -space-x-px text-sm"),
				gp.If(p.CurrentPage > 1, gh.Li(
					gh.A(gh.Href(prevURL.String()),
						gh.Class("flex items-center justify-center px-3 h-8 ms-0 leading-tight text-gray-500 bg-white border border-e-0 border-gray-300 rounded-s-lg hover:bg-gray-100 hover:text-gray-700"),
						gh.Span(gp.Text("Previous"))))),
				gp.Map(pages, func(i int) gp.Node {
					midValues := midURL.Query()
					midValues.Set("page", strconv.Itoa(i))
					midURL.RawQuery = midValues.Encode()
					if i != int(p.CurrentPage) {
						return gh.Li(gh.A(gp.Text(strconv.Itoa(i)), gh.Href(midURL.String()),
							gh.Class("flex items-center justify-center px-3 h-8 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700")))
					}
					return gh.Li(gh.A(gp.Text(strconv.Itoa(i)), gh.Href(midURL.String()),
						gh.Class("flex items-center justify-center px-3 h-8 text-blue-600 border border-gray-300 bg-blue-50 hover:bg-blue-100 hover:text-blue-700")))
				}),
				gp.If(p.CurrentPage < p.TotalPages, gh.Li(
					gh.A(gh.Href(nextURL.String()),
						gh.Class("flex items-center justify-center px-3 h-8 leading-tight text-gray-500 bg-white border border-gray-300 rounded-e-lg hover:bg-gray-100 hover:text-gray-700"),
						gh.Span(gp.Text("Next"))))),
			),
		)
	}
	summary := gh.Div(gh.Class("text-gray-700"),
		gp.Text("Showing "),
		gh.Span(gh.Class("font-medium"),
			gp.Text(strconv.Itoa(begin))),
		gp.Text(" to "),
		gh.Span(gh.Class("font-medium"),
			gp.Text(strconv.Itoa(end))),
		gp.Text(" of "),
		gh.Span(gh.Class("font-medium"),
			gp.Text(strconv.Itoa(numEntries))),
		gp.Text(" entries"),
	)
	concat := gh.Table(gh.Tr(gh.Td(summary), gh.Tr(gh.Td(paginator))))

	return &concat, nil

}

type ListViewer interface {
	List() gp.Node
}

type TableRow struct {
	Values []gp.Node
}

func (r TableRow) Render() gp.Node {
	return gh.Tr(gh.Class("text-neutral-1000"),
		gp.Map(r.Values, func(node gp.Node) gp.Node {
			return gh.Td(gh.Class("px-5 py-4 text-sm whitespace-nowrap"),
				node)
		}))

}

func MyTable(fields []string, rows []TableRow) gp.Node {
	return gh.Table(gh.Class("min-w-full divide-y divide-neutral-200"),
		TableHeader(fields),
		TableBody(rows),
	)
}

func TableHeader(fields []string) gp.Node {
	return gh.THead(gh.Tr(gh.Class("text-neutral-500"),
		gp.Map(fields, func(f string) gp.Node {
			return gh.Th(gh.Class("px-5 py-3 text-xs font-medium text-left uppercase"), gp.Text(f))
		})))
}

func TableBody(rows []TableRow) gp.Node {
	return gh.TBody(gh.Class("divide-y divide-neutral-200"),
		gp.Map(rows, func(r TableRow) gp.Node {
			return r.Render()

		}),
	)
}

func MakeDescriptionTable(keyvalues map[string]string, keys []string) gp.Node {
	var nodes []gp.Node
	for _, k := range keys {
		nodes = append(nodes, gh.Tr(gh.Class("text-neutral-500"),
			gh.Td(gh.Class("px-5 py-3 text-xs font-bold text-left uppercase border-b border-gray-300"), gp.Text(k)),
			gh.Td(gh.Class("w-30 border-b border-gray-300"), gp.Text(keyvalues[k]))))
	}
	table := gh.Div(gh.Class("max-w-xl w-full"),
		gh.Div(gh.Class("overflow-x-auto"),
			gh.Table(gh.Class("w-full border-collapse border border-gray-300"), gp.Group(nodes))))
	return table
}
