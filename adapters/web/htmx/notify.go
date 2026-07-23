package htmx

import (
	"encoding/json"
	"net/http"
)

func NotifySuccessPayloadAndReload(w http.ResponseWriter, title, message string) {
	payload, _ := json.Marshal(map[string]any{
		"htmx-notify-and-reload": map[string]string{
			"variant": "success",
			"title":   title,
			"message": message,
		},
	})
	w.Header().Set("HX-Trigger", string(payload))
	w.WriteHeader(http.StatusOK)
}

func NotifySuccessPayloadAndRedirect(w http.ResponseWriter, title, message, redirectURL string) {
	payload, _ := json.Marshal(map[string]any{
		"htmx-notify-and-redirect": map[string]string{
			"variant":  "success",
			"title":    title,
			"message":  message,
			"redirect": redirectURL,
		},
	})
	w.Header().Set("HX-Trigger", string(payload))
}

func NotifyError(w http.ResponseWriter, title, message string) {
	payload, _ := json.Marshal(map[string]any{
		"htmx-notify": map[string]string{
			"variant": "danger",
			"title":   title,
			"message": message,
		},
	})

	w.Header().Set("HX-Trigger", string(payload))
	w.WriteHeader(http.StatusUnprocessableEntity)
}
