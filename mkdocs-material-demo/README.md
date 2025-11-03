# MkDocs Material Demo

A simple demonstration of using the MkDocs Material module to build documentation sites.

## Usage

Build the example site and export it to a local directory:

```bash
dagger call --mod ./mkdocs-material-demo build-site export --path=./output
```

This will:

1. Load the fixture site from `fixtures/mkdocs-material/`
2. Build it using MkDocs Material
3. Export the static site to `./output`

## Example Output

The built site includes:

- `index.html` - The homepage
- `assets/` - CSS, JavaScript, and other assets
- `search/` - Search functionality
- `sitemap.xml` - Site map for search engines
