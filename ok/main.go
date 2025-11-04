// A Dagger module for ok (oslokommune/ok tool)

package main

import (
	"dagger/ok/internal/dagger"
)

type Ok struct{}

// Base returns a container with ok installed via mise
func (m *Ok) Base() *dagger.Container {
	return dag.Mise().Base().
		WithExec([]string{"mise", "use", "--global", "ubi:oslokommune/ok"}).
		WithEnvVariable("PATH", "/root/.local/share/mise/shims:$PATH", dagger.ContainerWithEnvVariableOpts{
			Expand: true,
		})
}
