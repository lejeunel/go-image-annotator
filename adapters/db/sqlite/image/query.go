package image

import (
	"errors"
	"strings"

	"go.tomakado.io/dumbql/query"
	"go.tomakado.io/dumbql/schema"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	qu "github.com/lejeunel/go-image-annotator/modules/query"
)

type windowExprs struct {
	prevId         string
	prevCollection string
	nextId         string
	nextCollection string
}

func makeAllWindowExprs(parser OrderStrParser, o im.OrderStr) (*windowExprs, error) {
	var (
		exprs windowExprs
		errs  []error
	)
	var err error
	exprs.prevId, err = makeWindowExpr(parser, "LAG(image_id, 1)", o, "prev_id")
	if err != nil {
		errs = append(errs, err)
	}
	exprs.prevCollection, err = makeWindowExpr(parser, "LAG(collection, 1)", o, "prev_collection")
	if err != nil {
		errs = append(errs, err)
	}
	exprs.nextId, err = makeWindowExpr(parser, "LEAD(image_id, 1)", o, "next_id")
	if err != nil {
		errs = append(errs, err)
	}
	exprs.nextCollection, err = makeWindowExpr(parser, "LEAD(collection, 1)", o, "next_collection")
	if err != nil {
		errs = append(errs, err)
	}
	return &exprs, errors.Join(errs...)

}

func makeWindowExpr(parser OrderStrParser, function string, ordering im.OrderStr, outName string) (string, error) {
	if ordering == "" {
		ordering = "image_id"
	}
	args, err := parser.Parse(ordering)
	if err != nil {
		return "", err
	}

	var terms []string
	for _, a := range args {
		if a.Order == im.DescOrder {
			terms = append(terms, a.Field+" "+"DESC")
		} else {
			terms = append(terms, a.Field)
		}
	}

	res := function + " OVER (ORDER BY " + strings.Join(terms, ",") + ") " + outName
	return res, nil
}

func MakeQueryParsers() (qu.FilterParser, qu.OrderParser) {
	sb := schema.NewSchemaBuilder()
	sb.AddField("collection", schema.Is[string]())
	sb.AddField("ingested_at", schema.Is[string]())
	sb.AddRegExpField(`^meta\..*$`, schema.Any(schema.Is[float64](), schema.Is[string](), schema.Is[bool]()))

	rb := query.NewRenamerBuilder()
	rb.Add(`\bmeta\.(.*)\b`, `json_extract(m.meta, '$.$1')`)

	filterParser := qu.NewFilterParser(
		sb.Build(),
		qu.WithRenamer(rb.Build()),
	)
	ob := qu.NewOrderParserBuilder()
	ob.AddField("image_id")
	ob.AddField("ingested_at")
	ob.AddField("collection")
	ob.AddRegExpField(`^meta\..*$`)
	ob.AddRenameRule(`\bmeta\.(.*)\b`, `json_extract(m.meta, '$.$1')`)

	orderParser := ob.Build()

	return filterParser, orderParser
}

type OrderStrParser interface {
	Parse(string) (im.OrderingArgs, error)
}

type FilterParser interface {
	ParseToSql(q string) (*qu.SQLizer, error)
}
