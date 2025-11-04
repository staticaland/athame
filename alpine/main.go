// A Dagger module for Alpine Linux containers

package main

import (
	"fmt"

	"dagger/alpine/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=alpine
	// +default="3.22.2@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412"
	imageTag string,
) *Alpine {
	return &Alpine{
		ImageTag: imageTag,
	}
}

type Alpine struct {
	ImageTag string
}

// Base returns a base Alpine Linux container
func (m *Alpine) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("alpine:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// WithPackages returns a container with the specified packages installed
func (m *Alpine) WithPackages(packages []string) *dagger.Container {
	args := append([]string{"apk", "add", "--no-cache"}, packages...)
	return m.Base().WithExec(args)
}
