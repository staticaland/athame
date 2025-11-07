// A Dagger module for triggering Render deploy hooks
//
// This module provides functions to trigger on-demand deploys of Render services
// via their deploy hook URLs. Deploy hooks are useful for CI/CD pipelines,
// image registry updates, and headless CMS integrations.

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"dagger/render-deploy-hook/internal/dagger"
)

// RenderDeployHook provides functions to trigger Render deploy hooks
type RenderDeployHook struct {
	// DeployHookURL is the unique, secret deploy hook URL for your Render service
	DeployHookURL *dagger.Secret
}

// New creates a new RenderDeployHook module instance
func New(
	// The unique, secret deploy hook URL for your Render service
	deployHookURL *dagger.Secret,
) *RenderDeployHook {
	return &RenderDeployHook{
		DeployHookURL: deployHookURL,
	}
}

// Deploy triggers a deploy of the Render service using its deploy hook URL.
// Returns the HTTP status code and response body.
func (m *RenderDeployHook) Deploy(ctx context.Context) (string, error) {
	hookURL, err := m.DeployHookURL.Plaintext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get deploy hook URL: %w", err)
	}
	return m.triggerDeploy(ctx, hookURL)
}

// DeployImage triggers a deploy of the Render service with a specific Docker image.
// This is useful for image-backed services where you want to deploy a specific
// image tag or digest.
//
// The imageURL should be in the format: registry/repository:tag or registry/repository@digest
// Example: docker.io/library/nginx:1.26 or docker.io/library/nginx@sha256:abc123...
func (m *RenderDeployHook) DeployImage(
	ctx context.Context,
	// The Docker image URL to deploy (e.g., docker.io/library/nginx:1.26)
	imageURL string,
) (string, error) {
	hookURL, err := m.DeployHookURL.Plaintext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get deploy hook URL: %w", err)
	}

	// URL encode the image URL for the imgURL query parameter
	encodedImageURL := url.QueryEscape(imageURL)

	// Append the imgURL query parameter
	fullURL := fmt.Sprintf("%s&imgURL=%s", hookURL, encodedImageURL)

	return m.triggerDeploy(ctx, fullURL)
}

// triggerDeploy sends an HTTP POST request to the deploy hook URL
func (m *RenderDeployHook) triggerDeploy(ctx context.Context, hookURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hookURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to trigger deploy: %w", err)
	}
	defer resp.Body.Close()

	return fmt.Sprintf("Deploy triggered successfully. Status: %d %s", resp.StatusCode, resp.Status), nil
}
