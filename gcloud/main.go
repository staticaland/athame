// A Dagger module for Google Cloud SDK (gcloud) commands
//
// This module provides a container-based interface to the gcloud CLI,
// allowing you to run Google Cloud commands in your Dagger pipelines.

package main

import (
	"context"
	"fmt"

	"dagger/gcloud/internal/dagger"
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
		"--region=" + region,
	}

	if allowUnauthenticated {
		args = append(args, "--allow-unauthenticated")
	}

	return ctr.WithExec(args).Stdout(ctx)
}
