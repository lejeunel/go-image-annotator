package query

import (
	"regexp"
)

type FieldName = string

type Field struct {
	Name        FieldName
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
