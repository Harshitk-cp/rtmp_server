package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type WebhookManager struct {
	webhooks map[string]string
	mu       sync.RWMutex
}

func NewWebhookManager() *WebhookManager {
	return &WebhookManager{
		webhooks: make(map[string]string),
	}
}

func (wm *WebhookManager) RegisterWebhook(streamKey, webhookURL string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.webhooks[streamKey] = webhookURL
	logrus.Info("Webhook registered successfully for stream key: ", streamKey)
}

func (wm *WebhookManager) UnregisterWebhook(streamKey string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	delete(wm.webhooks, streamKey)
}

func (wm *WebhookManager) SendWebhook(offlineStreamKey, event, ingressId string) {
	wm.mu.RLock()
	webhooks := make(map[string]string, len(wm.webhooks))
	for k, v := range wm.webhooks {
		webhooks[k] = v
	}
	wm.mu.RUnlock()

	payload := map[string]interface{}{
		"event": map[string]string{
			"offlineStreamKey": offlineStreamKey,
			"event":            event,
			"ingressId":        ingressId,
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Failed to marshal webhook payload: %v", err)
		return
	}

	for streamKey, webhookURL := range webhooks {
		go func(streamKey, webhookURL string) {
			logrus.Infof("Sending webhook to %s for stream key %s", webhookURL, streamKey)

			retries := 3
			for i := 0; i < retries; i++ {
				resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
				if err != nil {
					logrus.Errorf("Failed to send webhook (attempt %d): %v", i+1, err)
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					logrus.Infof("Webhook sent successfully to %s", webhookURL)
					return
				}
				logrus.Errorf("Webhook request failed with status code: %d (attempt %d)", resp.StatusCode, i+1)
			}
			logrus.Errorf("Failed to send webhook after %d attempts", retries)
		}(streamKey, webhookURL)
	}
}
