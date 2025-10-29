// A Dagger module for ORAS (OCI Registry As Storage)

package main

import (
	"fmt"

	"dagger/oras/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=ghcr.io/oras-project/oras
	// +default="v1.3.0@sha256:6ce045ce069a89934d6666b8b49f9c4c0145201bd6de6dbe2aee267814c55468"
	imageTag string,
) *Oras {
	return &Oras{
		ImageTag: imageTag,
	}
}

type Oras struct {
	ImageTag string
}

// Base returns the base container with ORAS installed
func (m *Oras) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("ghcr.io/oras-project/oras:%s", m.ImageTag)).
		WithoutEntrypoint()
}
