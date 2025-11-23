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

// base returns a container with dependencies installed and source code mounted
func (m *MieleCi) base() *dagger.Container {
	return dag.Node().Base().
		WithWorkdir("/app").
		// Copy package files first for better caching
		WithFile("/app/package-lock.json", m.Source.File("package-lock.json")).
		WithFile("/app/package.json", m.Source.File("package.json")).
		WithExec([]string{"npm", "ci", "--include=dev"}).
		// Copy application code after dependencies are installed
		WithDirectory("/app", m.Source)
}

// Verify runs all local validation steps: lint, build, test, and scan
// This phase requires no credentials and can run locally
func (m *MieleCi) Verify(ctx context.Context) (*dagger.Container, error) {
	m.notify(ctx, "Running verification steps...", dagger.NtfySendOpts{
		Title:    "Verify: Started",
		Priority: "default",
		Tags:     "mag",
	})

	// Step 1: Build
	m.notify(ctx, "Building application...", dagger.NtfySendOpts{
		Title:    "Verify: Build",
		Priority: "default",
		Tags:     "hammer_and_wrench",
	})

	builtSite := m.Build()

	// Step 2: Test
	m.notify(ctx, "Running tests...", dagger.NtfySendOpts{
		Title:    "Verify: Test",
		Priority: "default",
		Tags:     "test_tube",
	})

	testOutput, err := m.base().
		WithExec([]string{"npm", "test"}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("tests failed: %w", err)
	}
	fmt.Printf("Test output:\n%s\n", testOutput)

	// Step 3: Scan - Build container and scan it
	m.notify(ctx, "Scanning container for vulnerabilities...", dagger.NtfySendOpts{
		Title:    "Verify: Scan",
		Priority: "default",
		Tags:     "shield",
	})

	// Build a container for scanning (just linux/amd64 for faster verification)
	scanContainer := dag.Container(dagger.ContainerOpts{Platform: "linux/amd64"}).
		// renovate: datasource=docker depName=nginx
		From("nginx:1.27.5-alpine3.21@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10").
		WithDirectory("/usr/share/nginx/html", builtSite).
		WithExposedPort(80)

	scanResult, err := dag.Trivy().ScanContainer(ctx, scanContainer, "scan-target")
	if err != nil {
		m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
			Title:    "Verify: Scan Failed",
			Priority: "high",
			Tags:     "warning",
		})
		return nil, fmt.Errorf("trivy scan failed: %w", err)
	}
	fmt.Printf("Trivy scan results:\n%s\n", scanResult)

	m.notify(ctx, "Verification completed successfully.", dagger.NtfySendOpts{
		Title:    "Verify: Completed",
		Priority: "default",
		Tags:     "white_check_mark",
	})

	return scanContainer, nil
}

// Build builds the Vite application and returns the dist directory
func (m *MieleCi) Build() *dagger.Directory {
	buildContainer := m.base().
		WithEnvVariable("NODE_ENV", "production").
		WithExec([]string{"npm", "run", "build"})

	// Return the dist directory
	return buildContainer.Directory("/app/dist")
}

// Publish builds multi-platform containers and publishes them to GHCR
// This phase requires GHCR token for authentication
func (m *MieleCi) Publish(
	ctx context.Context,
	// The built application directory to package into a container
	builtSite *dagger.Directory,
	// GitHub token for GHCR authentication (get with: gh auth token)
	ghcrToken *dagger.Secret,
) (string, error) {
	m.notify(ctx, "Publishing to GHCR...", dagger.NtfySendOpts{
		Title:    "Publish: Started",
		Priority: "default",
		Tags:     "package",
	})

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

	// Publish to GHCR
	imageAddr := fmt.Sprintf("ghcr.io/%s/athame/%s:%s", m.GhcrUsername, m.ImageName, m.Tag)
	addr, err := dag.Container().
		WithRegistryAuth("ghcr.io", m.GhcrUsername, ghcrToken).
		Publish(ctx, imageAddr, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
			Title:    "Publish: Failed",
			Priority: "high",
			Tags:     "warning",
		})
		return "", fmt.Errorf("failed to publish to GHCR: %w", err)
	}

	m.notify(ctx,
		fmt.Sprintf("Published to GHCR.\n\n**Image:**\n```\n%s\n```", addr),
		dagger.NtfySendOpts{
			Title:    "Publish: Completed",
			Priority: "default",
			Tags:     "white_check_mark",
			Markdown: true,
		})

	return addr, nil
}

// DeployToCloud deploys the published container image to Fly.io
// This phase requires Fly.io credentials
func (m *MieleCi) DeployToCloud(
	ctx context.Context,
	// The published image address (e.g., ghcr.io/user/image:tag@sha256:...)
	imageAddr string,
	// Fly.io app name
	flyioApp string,
	// Fly.io API token
	flyioToken *dagger.Secret,
	// +optional
	// +default="arn"
	flyioRegion string,
) error {
	m.notify(ctx, "Deploying to Fly.io...", dagger.NtfySendOpts{
		Title:    "Deploy: Started",
		Priority: "default",
		Tags:     "rocket",
	})

	_, err := dag.Flyio().Deploy(ctx, flyioApp, imageAddr, flyioToken, dagger.FlyioDeployOpts{
		PrimaryRegion: flyioRegion,
		InternalPort:  80,
	})
	if err != nil {
		m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
			Title:    "Deploy: Failed",
			Priority: "high",
			Tags:     "warning",
		})
		return fmt.Errorf("fly.io deploy failed: %w", err)
	}

	flyioUrl := fmt.Sprintf("https://%s.fly.dev", flyioApp)
	m.notify(ctx,
		fmt.Sprintf("Deployed to Fly.io.\n\n**App:** %s", flyioApp),
		dagger.NtfySendOpts{
			Title:    "Deploy: Completed",
			Priority: "default",
			Tags:     "white_check_mark",
			Actions:  fmt.Sprintf("view, View Site, %s", flyioUrl),
			Markdown: true,
		})

	return nil
}

// Deploy runs all phases: Verify (lint/build/test/scan), Publish (GHCR), and DeployToCloud (Fly.io)
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
	m.notify(ctx, "Starting CI/CD pipeline...", dagger.NtfySendOpts{
		Title:    "Miele CI/CD Started",
		Priority: "default",
		Tags:     "hourglass_flowing_sand",
	})

	// Phase 1: Verify (no credentials needed)
	// We need to get the built site from Verify, not the scanned container
	_, err := m.Verify(ctx)
	if err != nil {
		return "", fmt.Errorf("verify phase failed: %w", err)
	}

	// Build fresh for publishing (Verify already validated it)
	builtSite := m.Build()

	// Phase 2: Publish (requires GHCR token)
	addr, err := m.Publish(ctx, builtSite, ghcrToken)
	if err != nil {
		return "", fmt.Errorf("publish phase failed: %w", err)
	}

	// Phase 3: DeployToCloud (requires Fly.io credentials)
	if err := m.DeployToCloud(ctx, addr, flyioApp, flyioToken, flyioRegion); err != nil {
		return addr, fmt.Errorf("deploy phase failed: %w", err)
	}

	return addr, nil
}
