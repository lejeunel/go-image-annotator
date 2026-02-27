package generic

import (
	"regexp"
)

var ResourceNameRegExp = regexp.MustCompile("^[a-z\\-_0-9]*$")
