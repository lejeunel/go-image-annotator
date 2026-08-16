package fake

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	qu "github.com/lejeunel/go-image-annotator/modules/query"
)

type QueryStrValidator struct {
	Err error
	Got string
}

func (c *QueryStrValidator) Validate(q string) error {
	c.Got = q
	if c.Err != nil {
		return c.Err
	}
	return nil
}

type FilterStrParser struct {
	Err  error
	sql  string
	args []any
}

func (p FilterStrParser) ParseToSql(q string) (*qu.SQLizer, error) {
	if p.Err != nil {
		return nil, p.Err
	}
	sqlizer := qu.NewSQLizer(p.sql, p.args)
	return &sqlizer, nil
}

type OrderStrParser struct {
	Err  error
	Args im.OrderingArgs
}

func (o *OrderStrParser) Parse(q string) (im.OrderingArgs, error) {
	if o.Err != nil {
		return im.OrderingArgs{}, o.Err
	}
	return o.Args, nil
}
