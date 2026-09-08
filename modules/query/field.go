package query

import (
	"regexp"
)

type Field struct {
	Name        string
	Description string
}

type RegExpField struct {
	RegExp      *regexp.Regexp
	Description string
}

type FieldOption func(*Field)

type RegExpFieldOption func(*RegExpField)

type FieldDescription struct {
	Name        string
	Description string
}
