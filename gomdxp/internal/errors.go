package docs

import (
	"errors"
)

var ErrPageListEmpty = errors.New("no documentation pages provided")
var ErrParsing = errors.New("error reading documentation file")
var ErrMetaDataParsing = errors.New("error parsing meta-data from yaml section")
