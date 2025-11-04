// A Dagger module for mise (polyglot runtime manager)

package main

import (
	"dagger/mise/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=archlinux/archlinux
	// +default="base-20251019.0.437072@sha256:4524236733437ff1f35531147aa444b32f674d9f328aebe06d3511be575c80a3"
	imageTag string,
) *Mise {
	return &Mise{
		ImageTag: imageTag,
	}
}

type Mise struct {
	ImageTag string
}

// Base returns a container with mise installed
func (m *Mise) Base() *dagger.Container {
	return dag.Archlinux(dagger.ArchlinuxOpts{
		ImageTag: m.ImageTag,
	}).WithPackages([]string{"mise"})
}
