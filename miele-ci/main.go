// A Dagger module for Miele delay-start app CI/CD: building and publishing the Vite application
//
// This module builds the Vite application, publishes it as a container image to GHCR,
// and deploys it to Fly.io.

package main

import (
	"context"
	"fmt"
	"time"

	"dagger/miele-ci/internal/dagger"
)

func New(
	// +defaultPath="/fixtures/miele-delay-start"
	source *dagger.Directory,
	// +default="miele"
	imageName string,
	// +default="latest"
	tag string,
	// +default="staticaland"
	ghcrUsername string,
) *MieleCi {
	return &MieleCi{
		Source:       source,
		ImageName:    imageName,
		Tag:          tag,
		GhcrUsername: ghcrUsername,
	}
}

type MieleCi struct {
	Source       *dagger.Directory
	ImageName    string
	Tag          string
	GhcrUsername string
}

// notify sends a notification via ntfy and logs any errors without failing
func (m *MieleCi) notify(ctx context.Context, message string, opts dagger.NtfySendOpts) {
	_, err := dag.Ntfy().Send(ctx, "athame", message, opts)
	if err != nil {
		fmt.Printf("Failed to send notification '%s': %v\n", opts.Title, err)
	}
}

// Build builds the Vite application and returns the dist directory
func (m *MieleCi) Build() *dagger.Directory {
	// Use Node module for base image (Alpine-based)
	buildContainer := dag.Node().Base().
		WithWorkdir("/app").
		WithEnvVariable("NODE_ENV", "production").
		// Copy package files first for better caching
		WithFile("/app/package-lock.json", m.Source.File("package-lock.json")).
		WithFile("/app/package.json", m.Source.File("package.json")).
		WithExec([]string{"npm", "ci", "--include=dev"}).
		// Copy application code after dependencies are installed
		WithDirectory("/app", m.Source).
		WithExec([]string{"npm", "run", "build"})

	// Return the dist directory
	return buildContainer.Directory("/app/dist")
}

// Publish builds the site and publishes it as a container image to GHCR
func (m *MieleCi) Publish(
	ctx context.Context,
	// GitHub token for GHCR authentication (get with: gh auth token)
	ghcrToken *dagger.Secret,
) (string, error) {
	// Build the site once
	builtSite := m.Build()

	// Platforms to build for
	platforms := []dagger.Platform{
		"linux/amd64",
		"linux/arm64",
	}

	// Create platform-specific variants
	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		ctr := dag.Container(dagger.ContainerOpts{Platform: platform}).
			// renovate: datasource=docker depName=nginx
			From("nginx:1.27.5-alpine3.21@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10").
			WithDirectory("/usr/share/nginx/html", builtSite).
			WithExposedPort(80).
			WithLabel("org.opencontainers.image.title", m.ImageName).
			WithLabel("org.opencontainers.image.version", m.Tag).
			WithLabel("org.opencontainers.image.created", time.Now().String()).
			WithLabel("org.opencontainers.image.source", "https://github.com/staticaland/athame")

		platformVariants = append(platformVariants, ctr)
	}

	// Scan the first platform variant with Trivy before publishing
	m.notify(ctx, "Scanning container for vulnerabilities...", dagger.NtfySendOpts{
		Title:    "Trivy Security Scan Started",
		Priority: "default",
		Tags:     "shield",
	})

	scanResult, err := dag.Trivy().ScanContainer(ctx, platformVariants[0], "scan-target")
	if err != nil {
		m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
			Title:    "Trivy Security Scan Failed",
			Priority: "high",
			Tags:     "warning",
		})
		return "", fmt.Errorf("trivy scan failed: %w", err)
	}
	fmt.Printf("Trivy scan results:\n%s\n", scanResult)

	m.notify(ctx, "Security scan completed successfully.", dagger.NtfySendOpts{
		Title:    "Trivy Security Scan Completed",
		Priority: "default",
		Tags:     "white_check_mark",
	})

	// Publish to GHCR
	imageAddr := fmt.Sprintf("ghcr.io/%s/athame/%s:%s", m.GhcrUsername, m.ImageName, m.Tag)
	addr, err := dag.Container().
		WithRegistryAuth("ghcr.io", m.GhcrUsername, ghcrToken).
		Publish(ctx, imageAddr, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		return "", fmt.Errorf("failed to publish to GHCR: %w", err)
	}

	return addr, nil
}

// Deploy builds, publishes, and deploys the application to Fly.io
func (m *MieleCi) Deploy(
	ctx context.Context,
	// GitHub token for GHCR authentication (get with: gh auth token)
	ghcrToken *dagger.Secret,
	// Fly.io app name
	flyioApp string,
	// Fly.io API token
	flyioToken *dagger.Secret,
	// +optional
	// +default="arn"
	flyioRegion string,
) (string, error) {
	// Send notification that deployment is starting
	m.notify(ctx, "Starting build...", dagger.NtfySendOpts{
		Title:    "Miele CI/CD Started",
		Priority: "default",
		Tags:     "hourglass_flowing_sand",
	})

	// Build and publish
	addr, err := m.Publish(ctx, ghcrToken)
	if err != nil {
		m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
			Title:    "Image Publishing Failed",
			Priority: "high",
			Tags:     "warning",
		})
		return "", err
	}

	// Send notification that image is published
	m.notify(ctx,
		fmt.Sprintf("Published to GHCR.\n\n**Image:**\n```\n%s\n```", addr),
		dagger.NtfySendOpts{
			Title:    "Image Publishing Completed",
			Priority: "default",
			Tags:     "white_check_mark",
			Markdown: true,
		})

	// Deploy to Fly.io
	_, err = dag.Flyio().Deploy(ctx, flyioApp, addr, flyioToken, dagger.FlyioDeployOpts{
		PrimaryRegion: flyioRegion,
		InternalPort:  80,
	})
	if err != nil {
		m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
			Title:    "Fly.io Deploy Failed",
			Priority: "high",
			Tags:     "warning",
		})
		return "", fmt.Errorf("fly.io deploy failed: %w", err)
	}

	// Send success notification
	flyioUrl := fmt.Sprintf("https://%s.fly.dev", flyioApp)
	m.notify(ctx,
		fmt.Sprintf("Deployed to Fly.io.\n\n**App:** %s", flyioApp),
		dagger.NtfySendOpts{
			Title:    "Fly.io Deploy Completed",
			Priority: "default",
			Tags:     "white_check_mark",
			Actions:  fmt.Sprintf("view, View Site, %s", flyioUrl),
			Markdown: true,
		})

	return addr, nil
}
