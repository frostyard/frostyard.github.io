# Frostyard Site Redesign

**Date:** 2026-02-25
**Status:** Approved

## Summary

Replace the current MkDocs + git submodules documentation site with a custom Go static site generator. All content lives in this repository as markdown files with YAML frontmatter, rendered through Templ templates, styled with Tailwind CSS, with Pagefind for client-side search.

## Goals

- Unified docs hub for the Frostyard ecosystem (no submodules)
- All content authored and managed in a single repository
- Full site: docs, getting started guides, blog, downloads, community, FAQ
- Frost/snow visual theme (tasteful, not campy)
- Static site deployed to GitHub Pages

## Tech Stack

- **Language:** Go
- **Templates:** Templ (github.com/a-h/templ)
- **Styling:** Tailwind CSS (standalone CLI, no Node)
- **Markdown:** goldmark with GFM extensions
- **Syntax highlighting:** Chroma
- **Frontmatter:** gopkg.in/yaml.v3
- **File watching:** fsnotify (dev server)
- **Search:** Pagefind (post-build CLI)
- **Deployment:** GitHub Pages (static files)

## Project Structure

```
site/
├── cmd/frostyard/main.go        # CLI: build, serve, new
├── internal/
│   ├── build/                   # Build orchestration, sitemap, RSS
│   ├── content/                 # Markdown + frontmatter parsing, content types
│   ├── render/                  # Templ rendering to static HTML
│   └── server/                  # Dev server with live reload
├── templates/
│   ├── layouts/                 # base, docs, blog, landing templ files
│   ├── components/              # nav, sidebar, search, footer, toc
│   └── pages/                   # home, downloads, community
├── content/
│   ├── docs/
│   │   ├── getting-started/     # Installation, quickstart
│   │   ├── images/              # Snow, Snowfield, Cayo, Snow Loaded
│   │   ├── tools/               # nbc, Chairlift, First Setup
│   │   └── faq.md
│   └── blog/posts/
├── static/                      # Images, fonts (copied to dist/)
├── tailwind.config.js
├── input.css
├── go.mod
├── Justfile
└── dist/                        # Generated output
```

## Content Model

### Page Frontmatter

```yaml
---
title: "Installing nbc"
description: "How to install nbc on your system"
section: tools/nbc
weight: 10
draft: false
---
```

### Section Index (`_index.md`)

```yaml
---
title: "Images"
description: "Frostyard bootc container images"
icon: "server"
---
```

### Blog Post

```yaml
---
title: "Introducing Snow 2.0"
date: 2026-02-20
author: "bjk"
tags: ["snow", "release"]
description: "What's new in Snow 2.0"
---
```

### Markdown Features

- GitHub Flavored Markdown (tables, task lists, strikethrough)
- Syntax highlighting with Chroma
- Auto-generated heading IDs
- Admonitions/callouts (info, warning, tip)

## Site Sections

- **Home** — Landing page explaining Frostyard
- **Docs**
  - **Getting Started** — Installation, quickstart
  - **Images** — Snow, Snowfield, Cayo, Snow Loaded
  - **Tools** — nbc, Chairlift, First Setup
  - **FAQ**
- **Blog** — Release announcements, guides, updates
- **Downloads** — Links to container images, ISOs, releases
- **Community** — Links, contribution guide

## Visual Design

- **Palette:** Cool blue-grays (slate-50 through slate-900) with ice blue accent (sky-400/sky-500)
- **Dark mode:** Default dark (slate-900 bg), light mode toggle
- **Typography:** System font stack for body, monospace for code
- **Frost touches:** Thin frost-blue gradient line at page top, cool-tinted code block backgrounds, faint blue-gray sidebar dividers
- **Navigation:** Fixed top bar (logo, sections, dark mode toggle, search). Docs pages have left sidebar (section tree) and right sidebar (page TOC)
- **Responsive:** Sidebars collapse on mobile, hamburger nav, TOC moves to top

## CLI Commands

### `frostyard build`

1. Glob `.md` files from `content/`
2. Parse frontmatter + markdown
3. Build section tree for navigation
4. Render through Templ templates to HTML
5. Run Tailwind CLI for CSS
6. Copy `static/` to `dist/`
7. Generate sitemap.xml and RSS feed
8. Run Pagefind for search index

### `frostyard serve`

1. Run full build
2. Serve `dist/` on localhost:3000
3. Watch content/, templates/, static/ for changes
4. Rebuild on change with live reload (WebSocket/SSE)

### `frostyard new <type> <path>`

- `frostyard new page docs/tools/nbc/install.md` — scaffold page
- `frostyard new post "My Blog Title"` — scaffold blog post

## Deployment

- GitHub Pages at `https://frostyard.github.io/`
- Justfile wraps build/serve/deploy
- GitHub Actions workflow for CI builds

## Migration

- Remove all git submodules and `.gitmodules`
- Remove MkDocs config, requirements.txt, .venv/
- Remove built site/ directory
- Remove submodule update workflow
- Port existing docs (nbc, chairlift, SNOW) into content/
- Keep git history

## Out of Scope

- i18n/localization
- Versioned docs
- Authentication / gated content
- CMS / admin interface
- Comments system
- Analytics (can add later)
