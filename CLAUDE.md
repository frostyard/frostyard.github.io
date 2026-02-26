# Frostyard Site

Custom Go static site generator. See `README.md` for content authoring and project layout.

## Commands

```bash
just build          # templ generate + build site to dist/
just serve          # templ generate + dev server on :3000
just test           # go test ./... -v
just generate       # templ generate only
just new-page PATH  # scaffold docs page
just new-post TITLE # scaffold blog post
```

## Architecture

Content flows: Markdown + YAML frontmatter → `internal/content` (parse) → `internal/render` (bridge to templates) → `templates/` (Templ components) → HTML in `dist/`.

- `internal/content/` — parser, loader, section tree builder. Two-pass: sections from `_index.md`, then pages assigned to sections.
- `internal/render/` — converts `content.Section`/`content.Page` to template-specific structs (`SidebarSection`, `TOCHeading`, etc.)
- `templates/components/` — reusable Templ components (Sidebar, Icon, TOC, Nav, Footer)
- `templates/layouts/` — page layouts (Base, Docs, Blog, Landing)
- `templates/pages/` — static pages (Home, Downloads, Community)

## Workflow

After editing `.templ` files, always run `templ generate` before `go build` or `go test`. The Justfile commands handle this automatically.

Templ generates `*_templ.go` files — these are committed to git alongside their `.templ` sources.

## Key Conventions

- URL paths always have leading and trailing slashes: `/docs/tools/nbc/`
- Sections need `_index.md` to appear in sidebar navigation
- Pages/sections sorted by `weight` (ascending, 0 = default, sorts last)
- Syntax highlighting uses Chroma with the `dracula` theme
- Tailwind CSS for all styling — no custom CSS classes
- Templ switch syntax uses colons (`case "x":`) not braces

## Testing

All tests are in `internal/content/`. Run with `go test ./... -v`. No tests exist yet for `render/`, `server/`, or `templates/`.

## Icons

Sidebar section icons use Heroicons (16x16 solid) rendered as inline SVG in `templates/components/icon.templ`. The `icon` frontmatter field in `_index.md` maps to icon names. See `/add-icon` skill for adding new icons.
