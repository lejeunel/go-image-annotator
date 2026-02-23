package tests

import (
	e "datahub/errors"
	g "datahub/generic"
	"testing"
)

func TestPagination(t *testing.T) {
	tests := map[string]struct {
		page        int64
		pageSize    int
		maxPageSize int
		success     bool
	}{
		"page inferior to 0 should fail":               {page: -3, pageSize: 2, maxPageSize: 2, success: false},
		"page inferior to 1 should fail":               {page: 0, pageSize: 2, maxPageSize: 2, success: false},
		"pagesize superior to max allowed should fail": {page: 1, pageSize: 3, maxPageSize: 2, success: false},
		"page and pagesize = 1 should succeed":         {page: 1, pageSize: 1, maxPageSize: 2, success: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := g.NewPaginationParams(tc.page, tc.pageSize, tc.maxPageSize)
			if tc.success == false {
				AssertErrorIs(t, err, e.ErrPagination)
			} else {
				AssertNoError(t, err)
			}
		})
	}

}
