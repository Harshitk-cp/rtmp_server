package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Harshitk-cp/rtmp_server/pkg/webhook"
)

func HandleRegisterWebhook(wm *webhook.WebhookManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			WebhookURL string `json:"webhookURL"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.WebhookURL == "" {
			http.Error(w, "WebhookURL is required", http.StatusBadRequest)
			return
		}

		wm.RegisterWebhook(req.WebhookURL)
		w.WriteHeader(http.StatusOK)
	}
}
