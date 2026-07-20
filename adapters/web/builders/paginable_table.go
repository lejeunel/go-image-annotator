package builders

import (
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	s "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type PaginableTableBuilder struct {
	listURL    string
	pagination s.Pagination
	tb.TableBuilder
}

func (b *PaginableTableBuilder) SetPagination(pg s.Pagination, url string) *PaginableTableBuilder {
	b.pagination = pg
	b.listURL = url
	return b
}

func (b *PaginableTableBuilder) AddRow(r tb.Row) *PaginableTableBuilder {
	b.TableBuilder.AddRow(r)
	return b
}

func NewPaginableTableBuilder(fields []string) PaginableTableBuilder {
	return PaginableTableBuilder{TableBuilder: tb.NewTableBuilder(fields)}
}
