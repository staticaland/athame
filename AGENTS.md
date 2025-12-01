# Dagger Development Guidelines

This is a Dagger monorepo where we create and manage multiple Dagger modules, always using the Go SDK.

## Quick Reference

- Always work from repository root
- Always use `--mod ./module-name` flag (except for repo module in `.dagger/`)
- Never use `cd` commands
- Import: `dagger/<module-name>/internal/dagger`
- Never edit: `dagger.gen.go`, `internal/`

## Directory Structure

- `.dagger/` - repo orchestration module (callable from root without `--mod`)
- `./module-name/` - individual modules

## Before You Start

Read the relevant doc BEFORE writing code:

- **Creating a new module**: Read `agent_docs/init-module.md`
- **Installing a module as a dependency**: Read `agent_docs/install-module.md`
- **Calling module functions or exporting results**: Read `agent_docs/call-and-export.md`
- **Sharing configuration across multiple functions**: Read `agent_docs/constructors.md`
- **Setting up a base container for a module**: Read `agent_docs/base-function.md`
- **Configuring container image versions**: Read `agent_docs/image-tags.md`
- **Making parameters optional or giving them default values**: Read `agent_docs/parameter-annotations.md`

## Finding Image Digests

Use the `container-image-lookup` subagent as your FIRST task when creating new container-based modules.

