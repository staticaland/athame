// A Dagger module for release-please - automated releases based on conventional commits

package main

import (
	"dagger/release-please/internal/dagger"
)

type ReleasePlease struct {
	// Source directory containing the repository
	Src *dagger.Directory
}

// New creates a new ReleasePlease module instance
func New(
	// Source directory containing the repository
	// +defaultPath="/"
	src *dagger.Directory,
) *ReleasePlease {
	return &ReleasePlease{
		Src: src,
	}
}

// Base returns the base container with Node.js and release-please installed
func (m *ReleasePlease) Base() *dagger.Container {
	return dag.Node(m.Src).
		Base().
		WithExec([]string{"npm", "install", "-g", "release-please"})
}

// Manifest runs release-please using manifest configuration files
// This creates pull requests and releases based on release-please-config.json
func (m *ReleasePlease) Manifest(
	// GitHub token for authentication
	token *dagger.Secret,
	// Repository in format "owner/repo"
	repoUrl string,
) *dagger.Container {
	return m.Base().
		WithSecretVariable("GITHUB_TOKEN", token).
		WithEnvVariable("REPO_URL", repoUrl).
		WithExec([]string{
			"release-please",
			"manifest-pr",
			"--token=$GITHUB_TOKEN",
			"--repo-url=" + repoUrl,
		}).
		WithExec([]string{
			"release-please",
			"manifest-release",
			"--token=$GITHUB_TOKEN",
			"--repo-url=" + repoUrl,
		})
}
