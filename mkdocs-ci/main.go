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

func New(
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
) *MkdocsCi {
	return &MkdocsCi{
		Source:       source,
		SitePath:     sitePath,
		ImageName:    imageName,
		Tag:          tag,
		GhcrUsername: ghcrUsername,
	}
}

type MkdocsCi struct {
	Source       *dagger.Directory
	SitePath     string
	ImageName    string
	Tag          string
	GhcrUsername string
}

// notify sends a notification via ntfy and logs any errors without failing
func (m *MkdocsCi) notify(ctx context.Context, message string, opts dagger.NtfySendOpts) {
	_, err := dag.Ntfy().Send(ctx, "athame", message, opts)
	if err != nil {
		fmt.Printf("Failed to send notification '%s': %v\n", opts.Title, err)
	}
}

// VerifyArtifact runs all local validation steps: lint, build, and scan
// Returns multi-platform container images ready for publishing
// This phase requires no credentials and can run locally
func (m *MkdocsCi) VerifyArtifact(
	ctx context.Context,
) ([]*dagger.Container, error) {
	siteDir := m.Source.Directory(m.SitePath)

	// Step 1: Lint - Run vale, prettier, markdownlint, and link checking concurrently
	m.notify(ctx, "Running linters...", dagger.NtfySendOpts{
		Title:    "Verify: Lint",
		Priority: "default",
		Tags:     "mag",
	})

	eg, gctx := errgroup.WithContext(ctx)

	// Run vale
	eg.Go(func() error {
		_, err := dag.Vale().Check(dagger.ValeCheckOpts{
			Source: siteDir,
			Path:   "docs",
		}).Stdout(gctx)
		return err
	})

	// Run prettier
	eg.Go(func() error {
		_, err := dag.Prettier().Check(dagger.PrettierCheckOpts{
			Source:  siteDir,
			Pattern: "docs/**/*.md",
		}).Stdout(gctx)
		return err
	})

	// Run markdownlint
	eg.Go(func() error {
		_, err := dag.MarkdownlintCli2().Check(dagger.MarkdownlintCli2CheckOpts{
			Source:  siteDir,
			Pattern: "docs/**/*.md",
		}).Stdout(gctx)
		return err
	})

	// Run link checking
	eg.Go(func() error {
		_, err := dag.Lychee().Check(dagger.LycheeCheckOpts{
			Source: siteDir,
			Path:   "docs",
		}).Stdout(gctx)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("lint failed: %w", err)
	}

	// Step 2: Build
	m.notify(ctx, "Building site...", dagger.NtfySendOpts{
		Title:    "Verify: Build",
		Priority: "default",
		Tags:     "hammer_and_wrench",
	})

	builtSite := m.Build()

	// Step 3: Build multi-platform containers
	m.notify(ctx, "Building container images...", dagger.NtfySendOpts{
		Title:    "Verify: Container Build",
		Priority: "default",
		Tags:     "package",
	})

	platforms := []dagger.Platform{
		"linux/amd64", // Required for Render and most cloud providers
		"linux/arm64", // For Apple Silicon Macs and ARM servers
	}

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

	// Step 4: Scan the first platform variant (amd64)
	m.notify(ctx, "Scanning container for vulnerabilities...", dagger.NtfySendOpts{
		Title:    "Verify: Scan",
		Priority: "default",
		Tags:     "shield",
	})

	scanResult, err := dag.Trivy().ScanContainer(ctx, platformVariants[0], "scan-target")
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

	// Return the built and scanned containers
	return platformVariants, nil
}

// deployToAllPlatforms deploys to Render, Fly.io, and Google Cloud Run concurrently based on provided credentials
func (m *MkdocsCi) deployToAllPlatforms(
	ctx context.Context,
	addr string,
	deployHookURL *dagger.Secret,
	flyioApp string,
	flyioToken *dagger.Secret,
	flyioRegion string,
	gcloudService string,
	gcloudProject string,
	gcloudRegion string,
	gcloudServiceAccountKey *dagger.Secret,
	gcloudAllowUnauthenticated bool,
	artifactRegistryRepo string,
	artifactRegistryRegion string,
) error {
	// Create error group for concurrent deployments
	eg, gctx := errgroup.WithContext(ctx)

	// Deploy to Render if provided
	if deployHookURL != nil {
		eg.Go(func() error {
			_, err := dag.RenderDeployHook(deployHookURL).Deploy(gctx)
			if err != nil {
				m.notify(gctx, "Check logs for details.", dagger.NtfySendOpts{
					Title:    "Render Deploy Failed",
					Priority: "high",
					Tags:     "warning",
				})
				return fmt.Errorf("render deploy failed: %w", err)
			}

			renderUrl := fmt.Sprintf("https://%s.onrender.com", m.ImageName)
			m.notify(gctx, "Deployed to Render.", dagger.NtfySendOpts{
				Title:    "Render Deploy Completed",
				Priority: "default",
				Tags:     "white_check_mark",
				Actions:  fmt.Sprintf("view, View Site, %s", renderUrl),
			})
			return nil
		})
	}

	// Deploy to Fly.io if provided
	if flyioApp != "" && flyioToken != nil {
		eg.Go(func() error {
			region := flyioRegion
			if region == "" {
				region = "arn" // default region
			}

			_, err := dag.Flyio().Deploy(gctx, flyioApp, addr, flyioToken, dagger.FlyioDeployOpts{
				PrimaryRegion: region,
				InternalPort:  80, // nginx default port
			})
			if err != nil {
				m.notify(gctx, "Check logs for details.", dagger.NtfySendOpts{
					Title:    "Fly.io Deploy Failed",
					Priority: "high",
					Tags:     "warning",
				})
				return fmt.Errorf("fly.io deploy failed: %w", err)
			}

			flyioUrl := fmt.Sprintf("https://%s.fly.dev", flyioApp)
			m.notify(gctx,
				fmt.Sprintf("Deployed to Fly.io.\n\n**App:** %s", flyioApp),
				dagger.NtfySendOpts{
					Title:    "Fly.io Deploy Completed",
					Priority: "default",
					Tags:     "white_check_mark",
					Actions:  fmt.Sprintf("view, View Site, %s", flyioUrl),
					Markdown: true,
				})
			return nil
		})
	}

	// Deploy to Google Cloud Run if provided
	if gcloudServiceAccountKey != nil && gcloudService != "" && gcloudProject != "" {
		eg.Go(func() error {
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

			_, err := dag.Gcloud().Deploy(gctx, gcloudService, artifactRegistryImage, gcloudProject, region, dagger.GcloudDeployOpts{
				AllowUnauthenticated: gcloudAllowUnauthenticated,
				ServiceAccountKey:    gcloudServiceAccountKey,
			})
			if err != nil {
				m.notify(gctx, "Check logs for details.", dagger.NtfySendOpts{
					Title:    "Google Cloud Run Deploy Failed",
					Priority: "high",
					Tags:     "warning",
				})
				return fmt.Errorf("google cloud run deploy failed: %w", err)
			}

			gcloudUrl := fmt.Sprintf("https://%s-%s.run.app", gcloudService, region)
			m.notify(gctx,
				fmt.Sprintf("Deployed to Cloud Run.\n\n**Service:** %s", gcloudService),
				dagger.NtfySendOpts{
					Title:    "Google Cloud Run Deploy Completed",
					Priority: "default",
					Tags:     "white_check_mark",
					Actions:  fmt.Sprintf("view, View Site, %s", gcloudUrl),
					Markdown: true,
				})
			return nil
		})
	}

	// Wait for all deployments to complete
	// Returns first error encountered, but all deployments run to completion
	return eg.Wait()
}

// Build builds the MkDocs Material site
func (m *MkdocsCi) Build() *dagger.Directory {
	docsSource := m.Source.Directory(m.SitePath)
	return dag.MkdocsMaterial().Build(dagger.MkdocsMaterialBuildOpts{
		Source: docsSource,
	})
}

// Publish runs VerifyArtifact, then publishes the verified containers to GHCR
// This phase requires GHCR token for authentication
func (m *MkdocsCi) Publish(
	ctx context.Context,
	// GitHub token for GHCR authentication (get with: gh auth token)
	ghcrToken *dagger.Secret,
) (string, error) {
	// Phase 1: VerifyArtifact (lint + build + scan) - returns built containers
	platformVariants, err := m.VerifyArtifact(ctx)
	if err != nil {
		return "", fmt.Errorf("verify artifact phase failed: %w", err)
	}

	// Phase 2: Publish the verified containers to GHCR
	m.notify(ctx, "Publishing to GHCR...", dagger.NtfySendOpts{
		Title:    "Publish: Started",
		Priority: "default",
		Tags:     "package",
	})

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
		fmt.Sprintf("Published to GHCR.\n\n**Image:**\n```\n%s\n```\n\n**Run:**\n```bash\ndocker run -p 8080:80 %s\n```", addr, addr),
		dagger.NtfySendOpts{
			Title:    "Publish: Completed",
			Priority: "default",
			Tags:     "white_check_mark",
			Markdown: true,
		})

	return addr, nil
}

// Deploy runs Publish, then deploys the container image to cloud platforms
// This phase requires GHCR token and cloud provider credentials
func (m *MkdocsCi) Deploy(
	ctx context.Context,
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
	m.notify(ctx, "Starting CI/CD pipeline...", dagger.NtfySendOpts{
		Title:    "MkDocs CI/CD Started",
		Priority: "default",
		Tags:     "hourglass_flowing_sand",
	})

	// Phase 1+2: Publish (which calls VerifyArtifact)
	addr, err := m.Publish(ctx, ghcrToken)
	if err != nil {
		return "", fmt.Errorf("publish phase failed: %w", err)
	}

	// Phase 3: Deploy to cloud platforms
	m.notify(ctx, "Deploying to cloud platforms...", dagger.NtfySendOpts{
		Title:    "Deploy: Started",
		Priority: "default",
		Tags:     "rocket",
	})

	err = m.deployToAllPlatforms(
		ctx,
		addr,
		deployHookURL,
		flyioApp,
		flyioToken,
		flyioRegion,
		gcloudService,
		gcloudProject,
		gcloudRegion,
		gcloudServiceAccountKey,
		gcloudAllowUnauthenticated,
		artifactRegistryRepo,
		artifactRegistryRegion,
	)

	if err != nil {
		m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
			Title:    "Deploy: Failed",
			Priority: "high",
			Tags:     "warning",
		})
		return addr, fmt.Errorf("deploy to cloud failed: %w", err)
	}

	m.notify(ctx, "All deployments completed successfully.", dagger.NtfySendOpts{
		Title:    "Deploy: Completed",
		Priority: "default",
		Tags:     "white_check_mark",
	})

	return addr, nil
}
