// A Dagger module for sending notifications via ntfy
package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Ntfy struct{}

// Send sends a notification to an ntfy topic
func (m *Ntfy) Send(
	ctx context.Context,
	// The ntfy topic to send to
	topic string,
	// The message to send
	message string,
	// +optional
	// The notification title
	title string,
	// +optional
	// +default="https://ntfy.sh"
	// The ntfy server URL
	server string,
	// +optional
	// Priority of the notification (e.g., "urgent", "high", "default", "low", "min")
	priority string,
	// +optional
	// Comma-separated list of tags (e.g., "warning,skull")
	tags string,
	// +optional
	// Enable Markdown formatting
	markdown bool,
	// +optional
	// Action buttons (format: "view, Label, URL" or "http, Label, URL")
	actions string,
) (string, error) {
	url := fmt.Sprintf("%s/%s", server, topic)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(message))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if title != "" {
		req.Header.Set("Title", title)
	}

	if priority != "" {
		req.Header.Set("Priority", priority)
	}

	if tags != "" {
		req.Header.Set("Tags", tags)
	}

	if markdown {
		req.Header.Set("Markdown", "yes")
	}

	if actions != "" {
		req.Header.Set("Actions", actions)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ntfy returned non-success status: %d", resp.StatusCode)
	}

	return fmt.Sprintf("Notification sent successfully to %s", topic), nil
}
