package components

import (
	"time"
)

func DateTimeToStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}
