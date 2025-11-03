// A Dagger module that demonstrates using MkDocs Material to build documentation sites
//
// This module shows how to use the mkdocs-material module to build static documentation sites.

package main

import (
	"dagger/mkdocs-material-demo/internal/dagger"
)

type MkdocsMaterialDemo struct{}

// BuildSite builds a MkDocs Material site from the provided source directory
func (m *MkdocsMaterialDemo) BuildSite(
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="fixtures/mkdocs-material"
	sitePath string,
) *dagger.Directory {
	// Navigate to the fixture directory
	docsSource := source.Directory(sitePath)

	// Use the mkdocs-material module to build the site
	return dag.MkdocsMaterial().Build(dagger.MkdocsMaterialBuildOpts{
		Source: docsSource,
	})
}
