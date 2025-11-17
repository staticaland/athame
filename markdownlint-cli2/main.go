// A Dagger module for markdownlint-cli2, a fast, flexible, configuration-based
// command-line interface for linting Markdown files with the markdownlint library

package main

import (
	"fmt"

	"dagger/markdownlint-cli-2/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=davidanson/markdownlint-cli2
	// +default="v0.18.1@sha256:173cb697a255a8a985f2c6a83b4f7a8b3c98f4fb382c71c45f1c52e4d4fed63a"
	imageTag string,
) *MarkdownlintCli2 {
	return &MarkdownlintCli2{
		ImageTag: imageTag,
	}
}

type MarkdownlintCli2 struct {
	ImageTag string
}

// Base returns the base container with markdownlint-cli2 installed
func (m *MarkdownlintCli2) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("davidanson/markdownlint-cli2:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// Check runs markdownlint-cli2 on the specified pattern
func (m *MarkdownlintCli2) Check(
	// +defaultPath="/"
	source *dagger.Directory,
	// Pattern to check (relative to source)
	// +default="docs/**/*.md"
	pattern string,
) *dagger.Container {
	return m.Base().
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"markdownlint-cli2", pattern})
}
