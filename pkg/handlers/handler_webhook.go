package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Harshitk-cp/rtmp_server/pkg/webhook"
)

func HandleRegisterWebhook(wm *webhook.WebhookManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			StreamKey  string `json:"streamKey"`
			WebhookURL string `json:"webhookURL"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		wm.RegisterWebhook(req.StreamKey, req.WebhookURL)
		w.WriteHeader(http.StatusOK)
	}
}
