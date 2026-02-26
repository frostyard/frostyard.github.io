# Frostyard Site

Custom static site generator for [frostyard.github.io](https://frostyard.github.io). Built with Go, [Templ](https://templ.guide), and Tailwind CSS.

## Prerequisites

- Go 1.24+
- [Templ](https://templ.guide) CLI (`go install github.com/a-h/templ/cmd/templ@latest`)
- [Tailwind CSS](https://tailwindcss.com) standalone CLI (download as `./tailwindcss`)
- [Pagefind](https://pagefind.app) (optional, for search indexing)
- [Just](https://just.systems) (optional, for task runner shortcuts)

## Quick Start

```bash
# Generate Templ code and build
just build

# Or without just:
templ generate
go run ./cmd/frostyard build

# Start dev server with live reload on :3000
just serve
```

The built site is output to `dist/`.

## Content Structure

All content lives in `content/` as Markdown files with YAML frontmatter. The directory structure directly determines the URL structure and sidebar navigation.

### How Directories Map to URLs

```
content/
  docs/
    _index.md            -> /docs/              (section index)
    faq.md               -> /docs/faq/
    getting-started/
      _index.md          -> /docs/getting-started/  (section index)
      quickstart.md      -> /docs/getting-started/quickstart/
    tools/
      _index.md          -> /docs/tools/
      nbc/
        _index.md        -> /docs/tools/nbc/
        encryption.md    -> /docs/tools/nbc/encryption/
    images/
      _index.md          -> /docs/images/
  blog/
    _index.md            -> /blog/
    posts/
      2025-01-15-hello.md -> /blog/posts/2025-01-15-hello/
```

Every `.md` file becomes a page at a URL derived from its path: strip `content/`, remove `.md`, add trailing slash.

### Sections and `_index.md`

Sections are directories that contain an `_index.md` file. The `_index.md` defines the section's title, description, and position in the sidebar. **A directory without `_index.md` will not appear in the navigation.**

```yaml
---
title: "Tools"
description: "Frostyard tools and utilities"
weight: 3
icon: "wrench"
---

Optional body text shown on the section's index page.
```

Sections can nest arbitrarily deep. Each level needs its own `_index.md`.

### Regular Pages

Any `.md` file that is not `_index.md` is a regular page. It belongs to the section defined by the nearest parent `_index.md`.

```yaml
---
title: "Quickstart"
description: "Get up and running with Frostyard in minutes"
weight: 1
---

Markdown content here.
```

### Frontmatter Fields

| Field         | Type     | Used in        | Description                                      |
|---------------|----------|----------------|--------------------------------------------------|
| `title`       | string   | all pages      | Page title (required)                            |
| `description` | string   | all pages      | Short description for meta tags and section lists |
| `weight`      | int      | docs           | Sort order within a section (lower = first)      |
| `draft`       | bool     | all pages      | If `true`, page is excluded from the build       |
| `icon`        | string   | `_index.md`    | Icon identifier for the section                  |
| `date`        | string   | blog posts     | Publication date (`YYYY-MM-DD`)                  |
| `author`      | string   | blog posts     | Author name                                      |
| `tags`        | []string | blog posts     | List of tags                                     |

### Ordering

Pages and sections within a section are sorted by `weight` (ascending). Pages with `weight: 0` (default) sort after pages with explicit weights.

### Blog Posts

Blog posts go in `content/blog/posts/`. Name them `YYYY-MM-DD-slug.md`. Posts are sorted by date (newest first) and rendered with the blog layout.

### Adding Content

Scaffold new content with the CLI:

```bash
# New docs page
just new-page docs/guides/setup

# New blog post
just new-post "My First Post"

# Or without just:
go run ./cmd/frostyard new page docs/guides/setup
go run ./cmd/frostyard new post "My First Post"
```

To add a new docs section:

1. Create the directory under `content/docs/`
2. Add an `_index.md` with `title`, `description`, and `weight`
3. Add pages in that directory

Example — adding a "Guides" section:

```
content/docs/guides/_index.md     # weight: 2 to place it after Getting Started
content/docs/guides/networking.md # weight: 1
content/docs/guides/storage.md    # weight: 2
```

## Project Layout

```
cmd/frostyard/         CLI entry point (build, serve, new)
internal/
  build/               Build pipeline (render, tailwind, sitemap, RSS, pagefind)
  content/             Markdown parser, content loader, section tree builder
  render/              Bridges content data to Templ templates
  server/              Dev server with file watching and SSE live reload
templates/
  layouts/             Base, Docs, Blog, Landing page layouts (Templ)
  components/          Nav, Sidebar, TOC, Footer components (Templ)
  pages/               Static pages: Home, Downloads, Community (Templ)
content/               Markdown content (docs, blog)
static/                Static assets copied to dist/ as-is
input.css              Tailwind CSS configuration
dist/                  Build output (gitignored)
```

## Build Pipeline

`go run ./cmd/frostyard build` runs these steps in order:

1. Load and parse all Markdown files from `content/`
2. Build section tree from `_index.md` files
3. Render each page to HTML using Templ templates
4. Render static pages (Home, Downloads, Community)
5. Copy `static/` assets to `dist/`
6. Run Tailwind CSS to generate `dist/css/style.css`
7. Generate `sitemap.xml`
8. Generate `blog/feed.xml` (RSS)
9. Run Pagefind to build the search index

## Deployment

Pushes to `main` trigger the GitHub Actions workflow (`.github/workflows/deploy.yml`) which builds and deploys to GitHub Pages.
