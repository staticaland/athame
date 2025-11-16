// A Dagger module for Google Cloud SDK (gcloud) commands
//
// This module provides a container-based interface to the gcloud CLI,
// allowing you to run Google Cloud commands in your Dagger pipelines.

package main

import (
	"context"
	"fmt"
	"strings"

	run "cloud.google.com/go/run/apiv2"
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"dagger/gcloud/internal/dagger"
	"google.golang.org/api/option"
)

func New(
	// renovate: datasource=docker depName=google/cloud-sdk
	// +default="546.0.0-alpine@sha256:cbc3420643b13a8b12950d03d2b0d31c4e522cd3d7438bc10bd741fb9947419c"
	imageTag string,
) *Gcloud {
	return &Gcloud{
		ImageTag: imageTag,
	}
}

type Gcloud struct {
	ImageTag string
}

// Base returns the base container with Google Cloud SDK installed
//
// Note: This uses the gcloud CLI container image. Alternatively, we could use
// the Google Cloud Go SDK (cloud.google.com/go) for programmatic access to GCP services.
// See also:
// - https://github.com/google-github-actions/deploy-cloudrun for inspiration
// - https://github.com/google-github-actions/auth for authentication patterns
func (m *Gcloud) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("google/cloud-sdk:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// Deploy deploys a container to Google Cloud Run
//
// Cloud Run is a serverless platform that automatically scales your containerized applications.
func (m *Gcloud) Deploy(
	ctx context.Context,
	// Service name
	service string,
	// Container image (e.g., LOCATION-docker.pkg.dev/PROJECT_ID/REPO/IMAGE:TAG)
	image string,
	// GCP project ID
	project string,
	// GCP region
	region string,
	// +optional
	// Allow unauthenticated access (makes the service publicly accessible without IAM authentication)
	allowUnauthenticated bool,
	// +optional
	// Service account key file (JSON) for authentication
	serviceAccountKey *dagger.Secret,
) (string, error) {
	ctr := m.Base()

	// If a service account key is provided, configure authentication
	if serviceAccountKey != nil {
		ctr = ctr.
			WithMountedSecret("/tmp/key.json", serviceAccountKey).
			WithExec([]string{"gcloud", "auth", "activate-service-account", "--key-file=/tmp/key.json"})
	}

	// Build the deploy command
	args := []string{
		"gcloud", "run", "deploy", service,
		"--image=" + image,
		"--project=" + project,
		"--region=" + region,
	}

	if allowUnauthenticated {
		args = append(args, "--allow-unauthenticated")
	}

	return ctr.WithExec(args).Stdout(ctx)
}

// GetServiceUrl retrieves the actual Cloud Run service URL using the Go SDK
//
// This function uses the Google Cloud Run API to get the actual service URL,
// which includes the project number automatically assigned by Cloud Run.
// The URL format is: https://[service]-[project-number].[region].run.app
func (m *Gcloud) GetServiceUrl(
	ctx context.Context,
	// Service name
	service string,
	// GCP project ID
	project string,
	// GCP region
	region string,
	// +optional
	// Service account key file (JSON) for authentication
	serviceAccountKey *dagger.Secret,
) (string, error) {
	var opts []option.ClientOption

	// If a service account key is provided, use it for authentication
	if serviceAccountKey != nil {
		// Export the secret to a temporary file in the container
		keyPath := "/tmp/gcp-key.json"
		keyContent, err := serviceAccountKey.Plaintext(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to read service account key: %w", err)
		}

		// Write key to a temporary location (in memory)
		// Note: In a real scenario, this would be handled more securely
		opts = append(opts, option.WithCredentialsJSON([]byte(keyContent)))
	}

	// Create Cloud Run client
	client, err := run.NewServicesClient(ctx, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to create Cloud Run client: %w", err)
	}
	defer client.Close()

	// Build the service name in the required format:
	// projects/{project}/locations/{location}/services/{service}
	serviceName := fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, service)

	// Get the service details
	req := &runpb.GetServiceRequest{
		Name: serviceName,
	}

	svc, err := client.GetService(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get service details: %w", err)
	}

	// Extract the URL from the service
	if svc.Uri == "" {
		return "", fmt.Errorf("service URL not found in response")
	}

	// Clean up the URL (remove any trailing slashes)
	url := strings.TrimSuffix(svc.Uri, "/")

	return url, nil
}
