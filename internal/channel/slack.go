package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SlackChannel struct {
	client *http.Client
}

func NewSlackChannel() *SlackChannel {
	return &SlackChannel{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SlackChannel) Send(ctx context.Context, event Event, cfg ChannelConfig) error {
	if cfg.WebhookURL == nil {
		return fmt.Errorf("slack - webhook_url not configured")
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("slack - parse payload: %w", err)
	}

	text := event.Payload
	if cfg.MessageTemplate != nil {
		rendered, err := renderTemplate(*cfg.MessageTemplate, payload)
		if err != nil {
			return fmt.Errorf("slack - render message: %w", err)
		}
		text = rendered
	}

	slackPayload := map[string]any{
		"text": text,
	}
	if cfg.SlackChannel != nil {
		slackPayload["channel"] = *cfg.SlackChannel
	}

	body, err := json.Marshal(slackPayload)
	if err != nil {
		return fmt.Errorf("slack - marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, *cfg.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("slack - build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("slack - send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack - non-2xx response: %d", resp.StatusCode)
	}

	return nil
}
