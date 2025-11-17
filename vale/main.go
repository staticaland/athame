// A Dagger module for Vale, a syntax-aware linter for prose

package main

import (
	"fmt"

	"dagger/vale/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=jdkato/vale
	// +default="v3.12.0@sha256:d5e8108bfd238192a82f303349b95ce39f605354843bc94811e24da1fe8f8ee0"
	imageTag string,
) *Vale {
	return &Vale{
		ImageTag: imageTag,
	}
}

type Vale struct {
	ImageTag string
}

// Base returns the base container with Vale installed
func (m *Vale) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("jdkato/vale:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// Check runs vale linter on the specified directory
func (m *Vale) Check(
	// +defaultPath="/"
	source *dagger.Directory,
	// Path to check (relative to source)
	// +default="docs"
	path string,
) *dagger.Container {
	return m.Base().
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"vale", path})
}
