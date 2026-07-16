package pagination

import (
	"fmt"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"math"
)

type Pagination struct {
	TotalRecords int64
	Page         int64
	PageSize     int
	TotalPages   int64
}

func New(page int64, pageSize int, totalRecords int64) Pagination {
	return Pagination{TotalRecords: totalRecords,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int64(math.Ceil(float64(totalRecords) / float64(pageSize)))}
}

func Validate(page int64, pageSize int) error {
	if page < 1 {
		return fmt.Errorf("validating page number: found non-positive value (%v): %w", page, e.ErrValidation)
	}
	return nil
}

type PaginationParams struct {
	PageSize int
	Page     int64
}
