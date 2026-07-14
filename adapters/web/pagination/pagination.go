package pagination

import (
	"net/http"
	"strconv"
)

func GetPageFromRequest(r *http.Request) int64 {
	pageStr := r.URL.Query().Get("page")

	if pageStr == "" {
		return 1
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}

	return int64(page)
}
