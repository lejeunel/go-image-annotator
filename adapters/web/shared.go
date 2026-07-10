package web

import (
	"encoding/json"
)

func NotifySuccessPayload(title, message string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"htmx-notify-and-reload": map[string]string{
			"variant": "success",
			"title":   title,
			"message": message,
		},
	})
}
