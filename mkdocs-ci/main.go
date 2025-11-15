// A Dagger module for MkDocs CI/CD: linting, building, and publishing documentation sites
//
// This module provides functions to lint markdown with vale, prettier, and markdownlint-cli2,
// build MkDocs Material sites, and publish the site as a container image to GitHub Container Registry.

package main

import (
	"context"
	"fmt"
	"time"

	"dagger/mkdocs-ci/internal/dagger"

	"golang.org/x/sync/errgroup"
)

type MkdocsCi struct{}

// Vale runs vale linter on markdown files
func (m *MkdocsCi) Vale(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
) (string, error) {
	return dag.Vale().Base().
		WithMountedDirectory("/src", source).
		WithWorkdir(fmt.Sprintf("/src/%s", sitePath)).
		WithExec([]string{"vale", "docs"}).
		Stdout(ctx)
}

// Prettier checks markdown formatting
func (m *MkdocsCi) Prettier(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
) (string, error) {
	return dag.Prettier().Base().
		WithMountedDirectory("/src", source).
		WithWorkdir(fmt.Sprintf("/src/%s", sitePath)).
		WithExec([]string{"prettier", "--check", "docs/**/*.md"}).
		Stdout(ctx)
}

// Markdownlint runs markdownlint-cli2 on markdown files
func (m *MkdocsCi) Markdownlint(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
) (string, error) {
	return dag.MarkdownlintCli2().Base().
		WithMountedDirectory("/src", source).
		WithWorkdir(fmt.Sprintf("/src/%s", sitePath)).
		WithExec([]string{"markdownlint-cli2", "docs/**/*.md"}).
		Stdout(ctx)
}

// CheckLinks validates links in markdown files using lychee
func (m *MkdocsCi) CheckLinks(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
) (string, error) {
	return dag.Lychee().Base().
		WithMountedDirectory("/src", source).
		WithWorkdir(fmt.Sprintf("/src/%s", sitePath)).
		WithExec([]string{"lychee", "--verbose", "docs"}).
		Stdout(ctx)
}

// RunAllTests runs vale, prettier, markdownlint, and link checking concurrently
func (m *MkdocsCi) RunAllTests(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
) error {
	// Create error group
	eg, gctx := errgroup.WithContext(ctx)

	// Run vale
	eg.Go(func() error {
		_, err := m.Vale(gctx, source, sitePath)
		return err
	})

	// Run prettier
	eg.Go(func() error {
		_, err := m.Prettier(gctx, source, sitePath)
		return err
	})

	// Run markdownlint
	eg.Go(func() error {
		_, err := m.Markdownlint(gctx, source, sitePath)
		return err
	})

	// Run link checking
	eg.Go(func() error {
		_, err := m.CheckLinks(gctx, source, sitePath)
		return err
	})

	// Wait for all tests to complete
	// If any test fails, the error will be returned
	return eg.Wait()
}

// Build builds the MkDocs Material site
func (m *MkdocsCi) Build(
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
) *dagger.Directory {
	docsSource := source.Directory(sitePath)
	return dag.MkdocsMaterial().Build(dagger.MkdocsMaterialBuildOpts{
		Source: docsSource,
	})
}

// Publish builds the site and publishes it as a container image to GHCR
func (m *MkdocsCi) Publish(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
	// +default="mkdocs-demo"
	imageName string,
	// +default="latest"
	tag string,
	// +default="staticaland"
	ghcrUsername string,
	// GitHub token for GHCR authentication (get with: gh auth token)
	ghcrToken *dagger.Secret,
) (string, error) {
	// Build the site once
	builtSite := m.Build(source, sitePath)

	// Platforms to build for: linux/amd64 (required for Render) and linux/arm64 (for Apple Silicon)
	platforms := []dagger.Platform{
		"linux/amd64", // Required for Render and most cloud providers
		"linux/arm64", // For Apple Silicon Macs and ARM servers
	}

	// Create platform-specific variants
	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		ctr := dag.Container(dagger.ContainerOpts{Platform: platform}).
			// renovate: datasource=docker depName=nginx
			From("nginx:1.27.5-alpine3.21@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10").
			WithDirectory("/usr/share/nginx/html", builtSite).
			WithExposedPort(80).
			WithLabel("org.opencontainers.image.title", imageName).
			WithLabel("org.opencontainers.image.version", tag).
			WithLabel("org.opencontainers.image.created", time.Now().String()).
			WithLabel("org.opencontainers.image.source", "https://github.com/staticaland/athame")

		platformVariants = append(platformVariants, ctr)
	}

	// Publish to GHCR
	imageAddr := fmt.Sprintf("ghcr.io/%s/athame/%s:%s", ghcrUsername, imageName, tag)
	addr, err := dag.Container().
		WithRegistryAuth("ghcr.io", ghcrUsername, ghcrToken).
		Publish(ctx, imageAddr, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		return "", fmt.Errorf("failed to publish to GHCR: %w", err)
	}

	return addr, nil
}

// DeployRender triggers a Render deploy hook after publishing
func (m *MkdocsCi) DeployRender(
	ctx context.Context,
	// The Render deploy hook URL (secret)
	deployHookURL *dagger.Secret,
) (string, error) {
	return dag.RenderDeployHook(deployHookURL).Deploy(ctx)
}

// DeployFlyio deploys the published image to Fly.io
func (m *MkdocsCi) DeployFlyio(
	ctx context.Context,
	// The fly.io app name
	app string,
	// The container image reference to deploy
	image string,
	// The fly.io API token
	token *dagger.Secret,
	// Primary region for the app (see https://fly.io/docs/reference/regions/)
	// +default="arn"
	primaryRegion string,
) (string, error) {
	return dag.Flyio().Deploy(ctx, app, image, token, dagger.FlyioDeployOpts{
		PrimaryRegion: primaryRegion,
		InternalPort:  80, // nginx default port
	})
}

// DeployGcloud deploys the published image to Google Cloud Run
func (m *MkdocsCi) DeployGcloud(
	ctx context.Context,
	// The Cloud Run service name
	service string,
	// The container image reference to deploy
	image string,
	// GCP project ID
	project string,
	// GCP region (see https://cloud.google.com/run/docs/locations)
	// +default="us-central1"
	region string,
	// Service account key file (JSON) for authentication
	// Example: file://$HOME/.config/gcloud/service-account-key.json
	// Or 1Password: cmd:op document get "Google Cloud - Service account key"
	serviceAccountKey *dagger.Secret,
	// Allow unauthenticated access (makes the service publicly accessible)
	// +default=true
	allowUnauthenticated bool,
) (string, error) {
	return dag.Gcloud().Deploy(ctx, service, image, project, region, dagger.GcloudDeployOpts{
		AllowUnauthenticated: allowUnauthenticated,
		ServiceAccountKey:    serviceAccountKey,
	})
}

// LintBuildPublish runs all tests concurrently, then builds and publishes if tests pass
func (m *MkdocsCi) LintBuildPublish(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
	// +default="mkdocs-demo"
	imageName string,
	// +default="latest"
	tag string,
	// +default="staticaland"
	ghcrUsername string,
	// GitHub token for GHCR authentication (get with: gh auth token)
	ghcrToken *dagger.Secret,
	// +optional
	deployHookURL *dagger.Secret,
	// +optional
	flyioApp string,
	// +optional
	flyioToken *dagger.Secret,
	// +optional
	flyioRegion string,
	// +optional
	gcloudService string,
	// +optional
	gcloudProject string,
	// +optional
	gcloudRegion string,
	// +optional
	gcloudServiceAccountKey *dagger.Secret,
	// +optional
	gcloudAllowUnauthenticated bool,
	// +optional
	// +default="ghcr"
	// Artifact Registry repository name
	artifactRegistryRepo string,
	// +optional
	// +default="europe-north2"
	// Artifact Registry region (can be different from Cloud Run region)
	artifactRegistryRegion string,
) (string, error) {
	// Send notification that deployment is starting
	_, err := dag.Ntfy().Send(
		ctx,
		"athame",
		"Starting MkDocs CI/CD pipeline - running tests...",
		dagger.NtfySendOpts{
			Title:    "Deployment Started",
			Priority: "default",
			Tags:     "rocket",
		},
	)
	if err != nil {
		// Log error but don't fail the pipeline
		fmt.Printf("Failed to send start notification: %v\n", err)
	}

	// Run all tests concurrently
	err = m.RunAllTests(ctx, source, sitePath)
	if err != nil {
		// Send failure notification
		_, notifyErr := dag.Ntfy().Send(
			ctx,
			"athame",
			fmt.Sprintf("Tests failed: %v", err),
			dagger.NtfySendOpts{
				Title:    "Tests Failed",
				Priority: "high",
				Tags:     "warning",
			},
		)
		if notifyErr != nil {
			fmt.Printf("Failed to send test failure notification: %v\n", notifyErr)
		}
		return "", fmt.Errorf("tests failed: %w", err)
	}

	// Send notification that tests passed
	_, err = dag.Ntfy().Send(
		ctx,
		"athame",
		"All tests passed! Building and deploying...",
		dagger.NtfySendOpts{
			Title:    "Tests Passed",
			Priority: "default",
			Tags:     "white_check_mark",
		},
	)
	if err != nil {
		fmt.Printf("Failed to send tests passed notification: %v\n", err)
	}

	// If tests pass, build and publish
	addr, err := m.Publish(ctx, source, sitePath, imageName, tag, ghcrUsername, ghcrToken)
	if err != nil {
		// Send deployment failure notification
		_, notifyErr := dag.Ntfy().Send(
			ctx,
			"athame",
			fmt.Sprintf("Deployment failed: %v", err),
			dagger.NtfySendOpts{
				Title:    "Deployment Failed",
				Priority: "urgent",
				Tags:     "warning",
			},
		)
		if notifyErr != nil {
			fmt.Printf("Failed to send deployment failure notification: %v\n", notifyErr)
		}
		return "", err
	}

	// Send notification that deployment is complete
	_, err = dag.Ntfy().Send(
		ctx,
		"athame",
		fmt.Sprintf("Deployment complete!\n\n**Image:**\n```\n%s\n```\n\n**Run locally:**\n```bash\ndocker run -p 8080:80 %s\n```", addr, addr),
		dagger.NtfySendOpts{
			Title:    "Deployment Complete",
			Priority: "default",
			Tags:     "white_check_mark",
			Markdown: true,
		},
	)
	if err != nil {
		fmt.Printf("Failed to send deployment complete notification: %v\n", err)
	}

	// Trigger Render deploy hook if provided
	if deployHookURL != nil {
		_, err := m.DeployRender(ctx, deployHookURL)
		if err != nil {
			// Send Render deploy failure notification
			_, notifyErr := dag.Ntfy().Send(
				ctx,
				"athame",
				fmt.Sprintf("Render deploy failed: %v", err),
				dagger.NtfySendOpts{
					Title:    "Render Deploy Failed",
					Priority: "high",
					Tags:     "warning",
				},
			)
			if notifyErr != nil {
				fmt.Printf("Failed to send Render deploy failure notification: %v\n", notifyErr)
			}
			return addr, fmt.Errorf("render deploy failed: %w", err)
		}

		// Send notification that Render deploy is complete
		renderUrl := fmt.Sprintf("https://%s.onrender.com", imageName)
		_, err = dag.Ntfy().Send(
			ctx,
			"athame",
			"Render deploy triggered successfully!",
			dagger.NtfySendOpts{
				Title:    "Render Deploy Complete",
				Priority: "default",
				Tags:     "rocket",
				Actions:  fmt.Sprintf("view, View Site, %s", renderUrl),
			},
		)
		if err != nil {
			fmt.Printf("Failed to send Render deploy complete notification: %v\n", err)
		}
	}

	// Deploy to Fly.io if provided
	if flyioApp != "" && flyioToken != nil {
		region := flyioRegion
		if region == "" {
			region = "arn" // default region
		}

		_, err := m.DeployFlyio(ctx, flyioApp, addr, flyioToken, region)
		if err != nil {
			// Send Fly.io deploy failure notification
			_, notifyErr := dag.Ntfy().Send(
				ctx,
				"athame",
				fmt.Sprintf("Fly.io deploy failed: %v", err),
				dagger.NtfySendOpts{
					Title:    "Fly.io Deploy Failed",
					Priority: "high",
					Tags:     "warning",
				},
			)
			if notifyErr != nil {
				fmt.Printf("Failed to send Fly.io deploy failure notification: %v\n", notifyErr)
			}
			return addr, fmt.Errorf("fly.io deploy failed: %w", err)
		}

		// Send notification that Fly.io deploy is complete
		flyioUrl := fmt.Sprintf("https://%s.fly.dev", flyioApp)
		_, err = dag.Ntfy().Send(
			ctx,
			"athame",
			fmt.Sprintf("Fly.io deploy successful!\nApp: %s", flyioApp),
			dagger.NtfySendOpts{
				Title:    "Fly.io Deploy Complete",
				Priority: "default",
				Tags:     "rocket",
				Actions:  fmt.Sprintf("view, View Site, %s", flyioUrl),
			},
		)
		if err != nil {
			fmt.Printf("Failed to send Fly.io deploy complete notification: %v\n", err)
		}
	}

	// Deploy to Google Cloud Run if service account key is provided
	if gcloudServiceAccountKey != nil && gcloudService != "" && gcloudProject != "" {
		region := gcloudRegion
		if region == "" {
			region = "us-central1" // default region
		}

		// Transform GHCR image to Artifact Registry remote repo format
		// From: ghcr.io/staticaland/athame/mkdocs-demo:latest@sha256:...
		// To: {artifactRegistryRegion}-docker.pkg.dev/{project}/{repo}/staticaland/athame/mkdocs-demo:latest@sha256:...
		// Extract the path after ghcr.io/
		ghcrPath := addr[len("ghcr.io/"):]
		artifactRegistryImage := fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s",
			artifactRegistryRegion, gcloudProject, artifactRegistryRepo, ghcrPath)

		_, err := m.DeployGcloud(ctx, gcloudService, artifactRegistryImage, gcloudProject, region, gcloudServiceAccountKey, gcloudAllowUnauthenticated)
		if err != nil {
			// Send Google Cloud deploy failure notification
			_, notifyErr := dag.Ntfy().Send(
				ctx,
				"athame",
				fmt.Sprintf("Google Cloud Run deploy failed: %v", err),
				dagger.NtfySendOpts{
					Title:    "GCloud Deploy Failed",
					Priority: "high",
					Tags:     "warning",
				},
			)
			if notifyErr != nil {
				fmt.Printf("Failed to send Google Cloud deploy failure notification: %v\n", notifyErr)
			}
			return addr, fmt.Errorf("google cloud run deploy failed: %w", err)
		}

		// Send notification that Google Cloud deploy is complete
		gcloudUrl := fmt.Sprintf("https://%s-%s.run.app", gcloudService, region)
		_, err = dag.Ntfy().Send(
			ctx,
			"athame",
			fmt.Sprintf("Google Cloud Run deploy successful!\n\n**Service:** %s\n**Region:** %s\n**Image:** %s", gcloudService, region, artifactRegistryImage),
			dagger.NtfySendOpts{
				Title:    "GCloud Deploy Complete",
				Priority: "default",
				Tags:     "rocket",
				Actions:  fmt.Sprintf("view, View Site, %s", gcloudUrl),
				Markdown: true,
			},
		)
		if err != nil {
			fmt.Printf("Failed to send Google Cloud deploy complete notification: %v\n", err)
		}
	}

	return addr, nil
}
