package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type WebhookManager struct {
	webhookURL string
	mu         sync.RWMutex
}

func NewWebhookManager() *WebhookManager {
	return &WebhookManager{}
}

func (wm *WebhookManager) RegisterWebhook(webhookURL string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.webhookURL = webhookURL
	logrus.Info("Webhook registered successfully: ", webhookURL)
}

func (wm *WebhookManager) SendWebhook(event, ingressId string) {
	wm.mu.RLock()
	webhookURL := wm.webhookURL
	wm.mu.RUnlock()

	if webhookURL == "" {
		logrus.Warn("No webhook URL registered")
		return
	}

	payload := map[string]interface{}{
		"event": map[string]string{
			"event":     event,
			"ingressId": ingressId,
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Failed to marshal webhook payload: %v", err)
		return
	}

	go func() {
		logrus.Infof("Sending webhook to %s", webhookURL)

		retries := 3
		for i := 0; i < retries; i++ {
			resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
			if err != nil {
				logrus.Errorf("Failed to send webhook (attempt %d): %v", i+1, err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				logrus.Infof("Webhook sent successfully to %s, for event %s", webhookURL, event)
				return
			}
			logrus.Errorf("Webhook request failed with status code: %d (attempt %d)", resp.StatusCode, i+1)
		}
		logrus.Errorf("Failed to send webhook after %d attempts", retries)
	}()
}
