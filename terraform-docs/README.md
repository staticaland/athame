# terraform-docs

Dagger module for generating Terraform documentation.

## Usage

```bash
# Generate markdown docs (default)
dagger call -m ./terraform-docs generate --source=./my-module export --path=.

# Generate JSON
dagger call -m ./terraform-docs generate --source=./my-module --format=json --output-file=docs.json export --path=.

# Process submodules recursively
dagger call -m ./terraform-docs generate --source=./my-module --recursive=true export --path=.
```

## Formats

Supports: `markdown`, `json`, `yaml`, `asciidoc`, `toml`, `xml`, `pretty`, `tfvars`

## Output Modes

- `inject` (default) - inserts between `<!-- BEGIN_TF_DOCS -->` markers
- `replace` - overwrites entire file
