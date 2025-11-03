// A Dagger module for MkDocs CI/CD: linting, building, and publishing documentation sites
//
// This module provides functions to lint markdown with vale, prettier, and markdownlint-cli2,
// build MkDocs Material sites, and publish the site as a container image to ttl.sh registry.

package main

import (
	"context"
	"fmt"

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

// RunAllTests runs vale, prettier, and markdownlint concurrently
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

// Publish builds the site and publishes it as a container image to ttl.sh
func (m *MkdocsCi) Publish(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
	// +default="mkdocs-demo"
	imageName string,
) (string, error) {
	// Build the site
	builtSite := m.Build(source, sitePath)

	// Create a container with nginx to serve the static site
	prodImage := dag.Container().
		// renovate: datasource=docker depName=nginx
		From("nginx:1.27.5-alpine3.21@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10").
		WithDirectory("/usr/share/nginx/html", builtSite).
		WithExposedPort(80)

	// Publish to ttl.sh registry
	addr, err := prodImage.Publish(ctx, fmt.Sprintf("ttl.sh/%s:1h", imageName))
	if err != nil {
		return "", err
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
) (string, error) {
	// Send notification that deployment is starting
	_, err := dag.Ntfy().Send(
		ctx,
		"athame",
		"Starting MkDocs CI/CD pipeline - running tests...",
		dagger.NtfySendOpts{
			Title:    "🚀 Deployment Started",
			Priority: "default",
			Tags:     "rocket,test",
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
				Title:    "❌ Tests Failed",
				Priority: "high",
				Tags:     "x,warning",
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
			Title:    "✅ Tests Passed",
			Priority: "default",
			Tags:     "white_check_mark,package",
		},
	)
	if err != nil {
		fmt.Printf("Failed to send tests passed notification: %v\n", err)
	}

	// If tests pass, build and publish
	addr, err := m.Publish(ctx, source, sitePath, imageName)
	if err != nil {
		// Send deployment failure notification
		_, notifyErr := dag.Ntfy().Send(
			ctx,
			"athame",
			fmt.Sprintf("Deployment failed: %v", err),
			dagger.NtfySendOpts{
				Title:    "❌ Deployment Failed",
				Priority: "urgent",
				Tags:     "x,warning",
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
			Title:    "🎉 Deployment Complete",
			Priority: "default",
			Tags:     "tada,white_check_mark",
			Markdown: true,
		},
	)
	if err != nil {
		fmt.Printf("Failed to send deployment complete notification: %v\n", err)
	}

	return addr, nil
}
