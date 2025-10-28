// A Dagger module for working with GitHub CLI
//
// This module provides functionality for using the GitHub CLI (gh)
// via asdf version manager.

package main

import (
	"context"

	"dagger/github-cli/internal/dagger"
)

func New(
	// renovate: datasource=github-releases depName=cli/cli
	// +default="2.71.0"
	version string,
) *GithubCli {
	return &GithubCli{
		Version: version,
	}
}

type GithubCli struct {
	Version string
}

// Base returns the base container with GitHub CLI installed via asdf
func (m *GithubCli) Base() *dagger.Container {
	return dag.Asdf().InstallPlugin(
		"github-cli",
		"https://github.com/bartlomiejdanek/asdf-github-cli.git",
		m.Version,
	)
}

// WithToken returns a container with the GitHub token configured
func (m *GithubCli) WithToken(token *dagger.Secret) *dagger.Container {
	return m.Base().
		WithSecretVariable("GITHUB_TOKEN", token)
}

// ListRepos returns a list of repositories for the authenticated user
func (m *GithubCli) ListRepos(
	ctx context.Context,
	token *dagger.Secret,
	// +optional
	// +default="100"
	limit string,
) (string, error) {
	return m.WithToken(token).
		WithExec([]string{"gh", "repo", "list", "--limit", limit}).
		Stdout(ctx)
}
