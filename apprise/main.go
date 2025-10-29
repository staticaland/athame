// A Dagger module for Apprise notification system
//
// Apprise allows you to send notifications to all of the most popular notification services.

package main

import (
	"fmt"

	"dagger/apprise/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=caronc/apprise
	// +default="1.2.2@sha256:0d74af8c1df9cf1de91f20f46d00ddee3a3efa15be179f9dbbe8a0f99d64268f"
	imageTag string,
) *Apprise {
	return &Apprise{
		ImageTag: imageTag,
	}
}

type Apprise struct {
	ImageTag string
}

// Base returns the base container with Apprise installed
func (m *Apprise) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("caronc/apprise:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// Send sends a notification via apprise
func (m *Apprise) Send(
	// Notification title
	title string,
	// Notification body/message
	body string,
	// Service URL to send notification to (e.g., discord://WEBHOOK_ID/WEBHOOK_TOKEN, syslog://)
	service *dagger.Secret,
) *dagger.Container {
	return m.Base().
		WithSecretVariable("APPRISE_SERVICE_URL", service).
		WithExec([]string{"sh", "-c", fmt.Sprintf("apprise -t '%s' -b '%s' \"$APPRISE_SERVICE_URL\"", title, body)})
}
