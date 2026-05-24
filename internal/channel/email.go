package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"
	tmpl "html/template"
	"strings"
)

type EmailChannel struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewEmailChannel(host, port, user, password, from string) *EmailChannel {
	return &EmailChannel{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
	}
}

func (e *EmailChannel) Send(ctx context.Context, event Event, cfg ChannelConfig) error {
	recipient, err := ExtractRecipient(event.Payload, cfg.ToField, cfg.ToFixed)
	if err != nil {
		return fmt.Errorf("email - extract recipient: %w", err)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("email - parse payload: %w", err)
	}

	subject, err := renderTemplate(subjectOrDefault(cfg.Subject), payload)
	if err != nil {
		return fmt.Errorf("email - render subject: %w", err)
	}

	body := ""
	if cfg.MessageTemplate != nil {
		rendered, err := renderTemplate(*cfg.MessageTemplate, payload)
		if err != nil {
			return fmt.Errorf("email - render body: %w", err)
		}
		body = rendered
	}

	msg := buildEmailMessage(e.from, recipient, subject, body)

	addr := fmt.Sprintf("%s:%s", e.host, e.port)
	auth := smtp.PlainAuth("", e.user, e.password, e.host)

	if err := smtp.SendMail(addr, auth, e.from, []string{recipient}, []byte(msg)); err != nil {
		return fmt.Errorf("email - send: %w", err)
	}

	return nil
}

func ExtractRecipient(payload string, toField *string, toFixed *string) (string, error) {
	if toField != nil && *toField != "" {
		var data map[string]any
		if err := json.Unmarshal([]byte(payload), &data); err != nil {
			return "", fmt.Errorf("parse payload json: %w", err)
		}
		parts := strings.Split(*toField, ".")
		val := traverseMap(data, parts)
		if val == "" {
			return "", fmt.Errorf("field %q not found in payload", *toField)
		}
		return val, nil
	}

	if toFixed != nil && *toFixed != "" {
		return *toFixed, nil
	}

	return "", fmt.Errorf("no recipient configured")
}

func traverseMap(data map[string]any, parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	val, ok := data[parts[0]]
	if !ok {
		return ""
	}
	if len(parts) == 1 {
		return fmt.Sprintf("%v", val)
	}
	nested, ok := val.(map[string]any)
	if !ok {
		return ""
	}
	return traverseMap(nested, parts[1:])
}

func renderTemplate(text string, data map[string]any) (string, error) {
	t, err := tmpl.New("t").Parse(text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func subjectOrDefault(subject *string) string {
	if subject != nil {
		return *subject
	}
	return "Notification"
}

func buildEmailMessage(from, to, subject, body string) string {
	return fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s", from, to, subject, body)
}
