package locations

import (
	"errors"
)

var ErrSiteNotFound = errors.New("site not found")
var ErrCameraNotFound = errors.New("camera not found")
