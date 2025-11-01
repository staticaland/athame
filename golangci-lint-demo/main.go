// A demo module for golangci-lint

package main

import (
	"context"

	"dagger/golangci-lint-demo/internal/dagger"
)

type GolangciLintDemo struct{}

// Lint runs golangci-lint on the provided source directory
func (m *GolangciLintDemo) Lint(
	ctx context.Context,
	// +defaultPath="/fixtures/hello-world-cli"
	source *dagger.Directory,
) (string, error) {
	return dag.GolangciLint().
		Base().
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run"}).
		Stdout(ctx)
}
