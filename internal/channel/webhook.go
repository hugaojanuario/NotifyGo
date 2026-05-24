package channel

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type WebhookChannel struct {
	client *http.Client
}

func NewWebhookChannel() *WebhookChannel {
	return &WebhookChannel{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WebhookChannel) Send(ctx context.Context, event Event, cfg ChannelConfig) error {
	if cfg.WebhookURL == nil {
		return fmt.Errorf("webhook - webhook_url not configured")
	}

	payload := []byte(event.Payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, *cfg.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("webhook - build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-NotifyGo-Topic", event.Topic)

	if cfg.WebhookSecret != nil && *cfg.WebhookSecret != "" {
		sig := computeHMAC(payload, *cfg.WebhookSecret)
		req.Header.Set("X-NotifyGo-Signature", "sha256="+sig)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook - send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook - non-2xx response: %d", resp.StatusCode)
	}

	return nil
}

func computeHMAC(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func generateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
