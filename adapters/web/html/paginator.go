package html

import (
	"fmt"
	"net/url"
	"strconv"

	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	s "github.com/lejeunel/go-image-annotator/shared"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakePaginatorEllipsis() Node {
	return Li(A(Href("#"),
		Class("flex size-6 items-center justify-center rounded-radius p-1 text-on-surface hover:text-primary dark:text-on-surface-dark dark:hover:text-primary-dark"),
		Aria("label", "more-pages"),
		Raw(`
				<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" aria-hidden="true" stroke="currentColor" class="size-6">
					<path stroke-linecap="round" stroke-linejoin="round" d="M6.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM12.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM18.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Z" />
				</svg>`),
	))
}

func MakePaginatorNumberedButton(baseURL url.URL, pageNumber int, isActive bool) Node {
	url := s.URLWithQuery(baseURL, "page", strconv.Itoa(pageNumber))
	class := Class("flex size-6 items-center justify-center rounded-radius p-1 text-on-surface hover:text-primary dark:text-on-surface-dark dark:hover:text-primary-dark")
	if isActive {
		class = Class("flex size-6 items-center justify-center rounded-radius bg-primary p-1 font-bold text-on-primary dark:bg-primary-dark dark:text-on-primary-dark")
	}
	return Li(A(Href(url.String()),
		class,
		Aria("label", fmt.Sprintf("page %v", pageNumber)),
		Text(strconv.Itoa(pageNumber)),
	))
}

func MakePaginator(baseURL string, currentPage, lastPage, numItems, totalItems int) Node {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		panic(fmt.Errorf("failed parsing pagination URL %v", baseURL))
	}
	prevURL := s.URLWithQuery(*parsedBaseURL, "page", strconv.Itoa(currentPage-1))
	nextURL := s.URLWithQuery(*parsedBaseURL, "page", strconv.Itoa(currentPage+1))
	return Nav(Aria("label", "pagination"),
		Ul(Class("flex shrink-0 items-center gap-2 text-sm font-medium"),
			If(currentPage > 1, Li(MakeNavigationButton(prevURL.String(), true, scr.ScrollPrevious, "Previous"))),
			If(currentPage > 1, MakePaginatorNumberedButton(*parsedBaseURL, currentPage-1, false)),
			If(lastPage > 1, MakePaginatorNumberedButton(*parsedBaseURL, currentPage, true)),
			If(lastPage-2 > currentPage, MakePaginatorEllipsis()),
			If(lastPage-1 > currentPage, MakePaginatorNumberedButton(*parsedBaseURL, lastPage-1, false)),
			If(lastPage > currentPage, MakePaginatorNumberedButton(*parsedBaseURL, lastPage, false)),
			If(currentPage < lastPage, Li(MakeNavigationButton(nextURL.String(), true, scr.ScrollNext, "Next"))),
			Span(Class("font-light"), Text(fmt.Sprintf("Showing %v item(s) out of %v", numItems, totalItems))),
		))
}
