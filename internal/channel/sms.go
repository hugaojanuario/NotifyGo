package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SMSChannel struct {
	accountSID string
	authToken  string
	from       string
	client     *http.Client
}

func NewSMSChannel(accountSID, authToken, from string) *SMSChannel {
	return &SMSChannel{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SMSChannel) Send(ctx context.Context, event Event, cfg ChannelConfig) error {
	recipient, err := ExtractRecipient(event.Payload, cfg.ToField, cfg.ToFixed)
	if err != nil {
		return fmt.Errorf("sms - extract recipient: %w", err)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("sms - parse payload: %w", err)
	}

	text := event.Payload
	if cfg.MessageTemplate != nil {
		rendered, err := renderTemplate(*cfg.MessageTemplate, payload)
		if err != nil {
			return fmt.Errorf("sms - render message: %w", err)
		}
		text = rendered
	}

	if len(text) > 160 {
		text = text[:160]
	}

	twilioURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.accountSID)

	data := url.Values{}
	data.Set("To", recipient)
	data.Set("From", s.from)
	data.Set("Body", text)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, twilioURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("sms - build request: %w", err)
	}
	req.SetBasicAuth(s.accountSID, s.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sms - send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sms - non-2xx response: %d", resp.StatusCode)
	}

	return nil
}
