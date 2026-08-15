package fake

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	ss "github.com/lejeunel/go-image-annotator/shared/sql"
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

type SQLizer struct {
	sql  string
	args []any
	Err  error
}

func (s SQLizer) ToSql() (string, []any, error) {
	return s.sql, s.args, s.Err
}

type FilterStrParser struct {
	Err error
	SQLizer
}

func (p *FilterStrParser) Parse(q string) (ss.SQLizer, error) {
	if p.Err != nil {
		return nil, p.Err
	}
	return p.SQLizer, nil
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
