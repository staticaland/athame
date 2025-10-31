// A Dagger module that demonstrates using LocalStack as a service
//
// This module shows how to use LocalStack as a service dependency in a pipeline.

package main

import (
	"context"

	"dagger/localstack-demo/internal/dagger"
)

type LocalstackDemo struct{}

// TestLocalstack starts LocalStack and makes HTTP requests to verify it's running
func (m *LocalstackDemo) TestLocalstack(ctx context.Context) (string, error) {
	// Start LocalStack service
	localstack := dag.Localstack().Run()

	// Create a container with curl and bind the LocalStack service
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithServiceBinding("localstack", localstack).
		WithExec([]string{"curl", "-s", "http://localstack:4566/_localstack/health"}).
		Stdout(ctx)
}

// CreateBucket demonstrates using AWS CLI with LocalStack to create an S3 bucket
func (m *LocalstackDemo) CreateBucket(
	ctx context.Context,
	// +default="demo-bucket"
	bucketName string,
) (string, error) {
	// Start LocalStack service
	localstack := dag.Localstack().Run()

	// Use AWS CLI configured for LocalStack to create S3 bucket
	return dag.AwsCli().
		LocalStack().
		WithServiceBinding("localstack", localstack).
		WithExec([]string{"aws", "s3", "mb", "s3://" + bucketName}).
		Stdout(ctx)
}

// TerraformApply demonstrates using Terraform with LocalStack to create infrastructure
// This example creates an S3 bucket using terraform-local
func (m *LocalstackDemo) TerraformApply(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/terraform-localstack"
	workdir string,
) (string, error) {
	// Start LocalStack service
	localstack := dag.Localstack().Run()

	// Use terraform-local to apply the Terraform configuration
	// Set S3_HOSTNAME to use the service binding hostname (enables path-style S3 access)
	return dag.Terraform().
		TerraformLocal().
		WithServiceBinding("localstack", localstack).
		WithEnvVariable("AWS_ENDPOINT_URL", "http://localstack:4566").
		WithEnvVariable("S3_HOSTNAME", "localstack").
		WithMountedDirectory("/work", source).
		WithWorkdir("/work/" + workdir).
		WithExec([]string{"tflocal", "init"}).
		WithExec([]string{"tflocal", "apply", "-auto-approve"}).
		Stdout(ctx)
}
