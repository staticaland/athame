# Calling Module Functions and Exporting Results

## Listing Functions

```bash
dagger functions --mod ./terraform-docs
```

## Calling Functions

```bash
dagger call --mod ./terraform-docs generate --source=.
```

## Exporting Results

```bash
dagger call --mod ./terraform-docs generate export --path=./output
```

## Directory Parameters

Use `.` for the repository root:

```bash
dagger call --mod ./terraform-docs generate --source=.
dagger call --mod ./mkdocs-ci build --docs=./docs
```

## Exploring Container Commands

Use `with-exec` with comma-separated `--args`:

```bash
dagger call --mod ./terraform-docs base with-exec --args="terraform-docs,--help" stdout
```

**Never use `terminal`** - it requires interactive TTY.

## Repo Module Exception

The repo module in `.dagger/` can be called without `--mod`:

```bash
dagger functions
dagger call terraform-docs export --path=./fixtures/terraform
```
