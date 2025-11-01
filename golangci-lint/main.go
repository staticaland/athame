// A Dagger module for golangci-lint

package main

import (
	"fmt"

	"dagger/golangci-lint/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=golangci/golangci-lint
	// +default="v2.6.0-alpine@sha256:1e8c410818ea9f1f4176b89dd2d95776f07184a7d4a8bf88d25e553b04c1995a"
	imageTag string,
) *GolangciLint {
	return &GolangciLint{
		ImageTag: imageTag,
	}
}

type GolangciLint struct {
	ImageTag string
}

// Base returns the base container with golangci-lint installed
func (m *GolangciLint) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golangci/golangci-lint:%s", m.ImageTag)).
		WithoutEntrypoint()
}
