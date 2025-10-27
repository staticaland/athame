// A Dagger module for Renovate

package main

import (
	"context"
	"fmt"

	"dagger/renovate/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=renovate/renovate
	// +default="41.163.0@sha256:0c1a0c9222430be38b2cf3136fec3b8c5ecf343807ee0026ee95e50db3e1ffb2"
	imageTag string,
) *Renovate {
	return &Renovate{
		ImageTag: imageTag,
	}
}

type Renovate struct {
	ImageTag string
}

// Base returns the base container with Renovate installed
func (m *Renovate) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("renovate/renovate:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// Run executes Renovate against a repository
func (m *Renovate) Run(
	ctx context.Context,
	// repository to run Renovate against (e.g., "owner/repo")
	project string,
	// authentication token for the platform
	token *dagger.Secret,
	// platform to use
	// +optional
	// +default="github"
	platform string,
) (string, error) {
	return m.Base().
		WithEnvVariable("RENOVATE_PLATFORM", platform).
		WithSecretVariable("RENOVATE_TOKEN", token).
		WithEnvVariable("RENOVATE_AUTODISCOVER", "false").
		WithEnvVariable("RENOVATE_REQUIRE_CONFIG", "optional").
		WithEnvVariable("LOG_LEVEL", "debug").
		WithExec([]string{"renovate", project}).
		Stdout(ctx)
}
