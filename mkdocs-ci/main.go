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
		"Starting tests...",
		dagger.NtfySendOpts{
			Title:    "MkDocs CI/CD Started",
			Priority: "default",
			Tags:     "hourglass_flowing_sand",
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
			"Check logs for details.",
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
		"Tests passed. Building site...",
		dagger.NtfySendOpts{
			Title:    "Tests Completed",
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
			"Check logs for details.",
			dagger.NtfySendOpts{
				Title:    "Image Publishing Failed",
				Priority: "high",
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
		fmt.Sprintf("Published to GHCR.\n\n**Image:**\n```\n%s\n```\n\n**Run:**\n```bash\ndocker run -p 8080:80 %s\n```", addr, addr),
		dagger.NtfySendOpts{
			Title:    "Image Publishing Completed",
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
		_, err := dag.RenderDeployHook(deployHookURL).Deploy(ctx)
		if err != nil {
			// Send Render deploy failure notification
			_, notifyErr := dag.Ntfy().Send(
				ctx,
				"athame",
				"Check logs for details.",
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
			"Deployed to Render.",
			dagger.NtfySendOpts{
				Title:    "Render Deploy Completed",
				Priority: "default",
				Tags:     "white_check_mark",
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

		_, err := dag.Flyio().Deploy(ctx, flyioApp, addr, flyioToken, dagger.FlyioDeployOpts{
			PrimaryRegion: region,
			InternalPort:  80, // nginx default port
		})
		if err != nil {
			// Send Fly.io deploy failure notification
			_, notifyErr := dag.Ntfy().Send(
				ctx,
				"athame",
				"Check logs for details.",
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
			fmt.Sprintf("Deployed to Fly.io.\n\n**App:** %s", flyioApp),
			dagger.NtfySendOpts{
				Title:    "Fly.io Deploy Completed",
				Priority: "default",
				Tags:     "white_check_mark",
				Actions:  fmt.Sprintf("view, View Site, %s", flyioUrl),
				Markdown: true,
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

		_, err := dag.Gcloud().Deploy(ctx, gcloudService, artifactRegistryImage, gcloudProject, region, dagger.GcloudDeployOpts{
			AllowUnauthenticated: gcloudAllowUnauthenticated,
			ServiceAccountKey:    gcloudServiceAccountKey,
		})
		if err != nil {
			// Send Google Cloud deploy failure notification
			_, notifyErr := dag.Ntfy().Send(
				ctx,
				"athame",
				"Check logs for details.",
				dagger.NtfySendOpts{
					Title:    "Google Cloud Run Deploy Failed",
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
			fmt.Sprintf("Deployed to Cloud Run.\n\n**Service:** %s", gcloudService),
			dagger.NtfySendOpts{
				Title:    "Google Cloud Run Deploy Completed",
				Priority: "default",
				Tags:     "white_check_mark",
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
