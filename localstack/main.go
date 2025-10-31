// A Dagger module for LocalStack
//
// LocalStack is a fully functional local AWS cloud stack that allows you to
// develop and test your cloud and serverless applications offline.

package main

import (
	"fmt"

	"dagger/localstack/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=localstack/localstack
	// +default="4.10.0@sha256:a65ee2a9d45a7a34a1f1faae515d2e577ce11210312c077700ccc82daefec238"
	imageTag string,
) *Localstack {
	return &Localstack{
		ImageTag: imageTag,
	}
}

type Localstack struct {
	ImageTag string
}

// Base returns the base container with LocalStack installed
func (m *Localstack) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("localstack/localstack:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// Run starts a LocalStack service
func (m *Localstack) Run() *dagger.Service {
	return m.Base().
		WithExposedPort(4566).
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{"docker-entrypoint.sh"},
		})
}
