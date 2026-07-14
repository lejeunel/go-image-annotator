package htmx

import (
	"encoding/json"
)

func NotifySuccessPayloadAndReload(title, message string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"htmx-notify-and-reload": map[string]string{
			"variant": "success",
			"title":   title,
			"message": message,
		},
	})
}

func NotifyError(title, message string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"htmx-notify": map[string]string{
			"variant": "danger",
			"title":   title,
			"message": message,
		},
	})
}
