// A Dagger module for Go CI/CD: linting, building, and publishing Go applications
//
// This module provides functions to lint Go code with golangci-lint,
// build Go binaries, and publish container images to ttl.sh registry.

package main

import (
	"context"
	"fmt"

	"dagger/go-ci/internal/dagger"

	"golang.org/x/sync/errgroup"
)

func New(
	// renovate: datasource=docker depName=golang
	// +default="1.25.3-alpine3.22@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34"
	golangImageTag string,
) *GoCi {
	return &GoCi{
		GolangImageTag: golangImageTag,
	}
}

type GoCi struct {
	GolangImageTag string
}

// Base returns the base container with Go installed
func (m *GoCi) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golang:%s", m.GolangImageTag)).
		WithoutEntrypoint()
}

// Lint runs golangci-lint on the provided Go source code
func (m *GoCi) Lint(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
) (string, error) {
	return dag.GolangciLint().Base().
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "./..."}).
		Stdout(ctx)
}

// Gosec runs gosec security scanner on the provided Go source code
func (m *GoCi) Gosec(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
) (string, error) {
	return dag.Gosec().Base().
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"gosec", "./..."}).
		Stdout(ctx)
}

// RunAllTests runs linter and gosec concurrently
func (m *GoCi) RunAllTests(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
) error {
	// Create error group
	eg, gctx := errgroup.WithContext(ctx)

	// Run linter
	eg.Go(func() error {
		_, err := m.Lint(gctx, source)
		return err
	})

	// Run gosec
	eg.Go(func() error {
		_, err := m.Gosec(gctx, source)
		return err
	})

	// Wait for all tests to complete
	// If any test fails, the error will be returned
	return eg.Wait()
}

// Build builds a Go binary and publishes it as a container image to ttl.sh
func (m *GoCi) Build(
	ctx context.Context,
	// source code location
	// +defaultPath="/"
	source *dagger.Directory,
	// binary name
	// +default="app"
	binaryName string,
	// image name for ttl.sh
	// +default="myapp"
	imageName string,
) (string, error) {
	// Build the Go binary
	builder := m.Base().
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", binaryName})

	// Create production image on alpine base
	prodImage := dag.Container().
		From("alpine:latest").
		WithFile(fmt.Sprintf("/bin/%s", binaryName), builder.File(fmt.Sprintf("/src/%s", binaryName))).
		WithEntrypoint([]string{fmt.Sprintf("/bin/%s", binaryName)})

	// Publish to ttl.sh registry
	addr, err := prodImage.Publish(ctx, fmt.Sprintf("ttl.sh/%s:latest", imageName))
	if err != nil {
		return "", err
	}

	return addr, nil
}

// LintAndBuild runs linting first, then builds and publishes if linting passes
func (m *GoCi) LintAndBuild(
	ctx context.Context,
	// source code location
	// +defaultPath="/"
	source *dagger.Directory,
	// binary name
	// +default="app"
	binaryName string,
	// image name for ttl.sh
	// +default="myapp"
	imageName string,
) (string, error) {
	// Run linting first
	_, err := m.Lint(ctx, source)
	if err != nil {
		return "", fmt.Errorf("linting failed: %w", err)
	}

	// If linting passes, build and publish
	return m.Build(ctx, source, binaryName, imageName)
}
