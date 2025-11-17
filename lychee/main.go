// A Dagger module for lycheeverse/lychee - a fast, async link checker
package main

import (
	"fmt"

	"dagger/lychee/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=lycheeverse/lychee
	// +default="0.15.1-alpine@sha256:214ed75d61117c5dc39310b9da73bb9fae5333f6f6eb6891e861e79cda780268"
	imageTag string,
) *Lychee {
	return &Lychee{
		ImageTag: imageTag,
	}
}

type Lychee struct {
	ImageTag string
}

// Base returns the base container with lychee installed
func (m *Lychee) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("lycheeverse/lychee:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// Check runs lychee link checker on the specified directory
func (m *Lychee) Check(
	// +defaultPath="/"
	source *dagger.Directory,
	// Path to check (relative to source)
	// +default="docs"
	path string,
) *dagger.Container {
	return m.Base().
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"lychee", "--verbose", path})
}
