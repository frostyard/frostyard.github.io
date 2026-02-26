# Sidebar Icons Design

## Problem

Content front matter includes an `icon` field that is parsed and stored but never rendered. Section icons should appear in the sidebar navigation.

## Decisions

- **Icon source:** Heroicons (outline, 16x16), SVGs copied directly into a Templ component
- **Rendering approach:** A Templ `Icon(name string)` component maps icon name strings to inline SVG markup
- **Icon depth:** Icons render wherever the `icon` field is set; authors control depth by only adding icons to sections up to two levels deep
- **Missing icons:** Sections without an `icon` field render no icon and no placeholder
- **Layout:** Section titles use `flex items-center gap-1.5` to align icon and text

## Data Flow

```
content/_index.md (icon: "server")
  → content.Section.Icon (already parsed)
    → convertSection passes Icon to SidebarSection.Icon
      → sidebar.templ renders @Icon(section.Icon) next to title
```

## Changes

| File | Change |
|------|--------|
| `templates/components/icon.templ` | **New.** Icon component with name→SVG mapping |
| `templates/components/sidebar.templ` | Add `Icon string` to `SidebarSection`, render icon in section titles |
| `internal/render/render.go` | Pass `sec.Icon` through in `convertSection` |

## Initial Icons

- `server` — Heroicons outline server icon
- `wrench` — Heroicons outline wrench icon
