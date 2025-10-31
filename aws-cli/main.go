// AWS CLI module for Dagger
//
// Provides a containerized AWS CLI for interacting with AWS services.
// Includes a LocalStack() function for local development against LocalStack.

package main

import (
	"fmt"

	"dagger/aws-cli/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=amazon/aws-cli
	// +default="2.31.26@sha256:cf1851fa3162c35009b2dc6d2df2797e5b0e9723fe546f545c9fa34a3dc03477"
	imageTag string,
) *AwsCli {
	return &AwsCli{
		ImageTag: imageTag,
	}
}

type AwsCli struct {
	ImageTag string
}

// Base returns the base container with AWS CLI installed
func (m *AwsCli) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("amazon/aws-cli:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// LocalStack returns a container configured to use LocalStack
//
// Sets test credentials, region, and endpoint URL.
// User should bind the LocalStack service when calling this.
func (m *AwsCli) LocalStack() *dagger.Container {
	return m.Base().
		WithEnvVariable("AWS_ACCESS_KEY_ID", "test").
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", "test").
		WithEnvVariable("AWS_DEFAULT_REGION", "us-east-1").
		WithEnvVariable("AWS_ENDPOINT_URL", "http://localstack:4566")
}
