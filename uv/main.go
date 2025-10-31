// A Dagger module for uv, the Python package manager

package main

import (
	"fmt"

	"dagger/uv/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=ghcr.io/astral-sh/uv
	// +default="0.9.7-alpine3.22@sha256:ce2e7e691797f9bd2ee1b15fe59d272cb26d9662eda746e0fc1542c74a558064"
	imageTag string,
) *Uv {
	return &Uv{
		ImageTag: imageTag,
	}
}

type Uv struct {
	ImageTag string
}

// Base returns the base container with uv installed
func (m *Uv) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("ghcr.io/astral-sh/uv:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// ToolInstall installs a Python tool using uv
func (m *Uv) ToolInstall(
	// Tool name to install
	name string,
	// Version constraint (e.g., "==0.5.0", ">=1.0.0")
	// +optional
	version string,
) *dagger.Container {
	toolSpec := name
	if version != "" {
		toolSpec = fmt.Sprintf("%s%s", name, version)
	}

	return m.Base().
		WithExec([]string{"uv", "tool", "install", toolSpec})
}
