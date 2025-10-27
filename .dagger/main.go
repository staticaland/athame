// Repo orchestration module

package main

import (
	"context"
	"dagger/repo/internal/dagger"
)

func New(
	// The source directory of the repository
	// +defaultPath="/"
	src *dagger.Directory,
) *Repo {
	return &Repo{
		Src: src,
	}
}

type Repo struct {
	Src *dagger.Directory
}

// TerraformPlan runs terraform plan on the fixtures/terraform directory
func (m *Repo) TerraformPlan(ctx context.Context) (string, error) {
	return dag.Terraform().
		Base().
		WithMountedDirectory("/src", m.Src.Directory("fixtures/terraform")).
		WithWorkdir("/src").
		WithExec([]string{"terraform", "init"}).
		WithExec([]string{"terraform", "plan"}).
		Stdout(ctx)
}

// TerraformDocs generates terraform documentation for the fixtures/terraform directory
func (m *Repo) TerraformDocs() *dagger.Directory {
	return dag.TerraformDocs().
		Base().
		WithMountedDirectory("/src", m.Src.Directory("fixtures/terraform")).
		WithWorkdir("/src").
		WithExec([]string{"sh", "-c", "terraform-docs markdown . > README.md"}).
		Directory("/src")
}

// Boilerplate runs the boilerplate templating tool with a template source
func (m *Repo) Boilerplate(
	ctx context.Context,
	// Template source (git URL with optional version and subpath, e.g., https://github.com/username/repo.git#version:subpath)
	// +default="https://github.com/gruntwork-io/boilerplate.git#main:examples/for-learning-and-testing/terraform"
	templateSrc string,
	// Output folder path
	// +default="output"
	outputFolder string,
) *dagger.Directory {
	template := dag.Git(templateSrc).Head().Tree()

	return dag.Boilerplate().
		Base().
		WithMountedDirectory("/template", template).
		WithWorkdir("/work").
		WithExec([]string{"sh", "-c", "ls -la /template"}).
		WithExec([]string{"boilerplate", "--template-url", "/template", "--output-folder", outputFolder, "--var", "ServerName=MyServer", "--non-interactive"}).
		Directory("/work/" + outputFolder)
}
