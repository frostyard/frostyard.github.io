# Extensions Documentation Section Design

## Overview

Add an "Extensions" section to the Frostyard documentation site, populated from the systemd sysext definitions in the [snosi repository](https://github.com/frostyard/snosi/tree/main/mkosi.images).

## Source Data

The snosi repo defines 10 mkosi images. The `base` image is the foundation layer (not a sysext), so it is excluded. The remaining 9 are systemd system extensions:

1. **1password-cli** — 1Password CLI tool
2. **debdev** — Debian development/bootstrapping tools
3. **dev** — General development toolchain
4. **docker** — Docker CE with plugins
5. **emdash** — Emdash desktop application
6. **incus** — Incus container/VM manager
7. **nix** — Nix package manager
8. **podman** — Podman container runtime with ecosystem tools
9. **tailscale** — Tailscale VPN

## Structure

Flat section with individual pages (matches existing site patterns):

```
content/docs/extensions/
├── _index.md          (section intro, weight: 3, icon: puzzle-piece)
├── podman.md          (weight: 1)
├── docker.md          (weight: 2)
├── incus.md           (weight: 3)
├── dev.md             (weight: 4)
├── debdev.md          (weight: 5)
├── nix.md             (weight: 6)
├── tailscale.md       (weight: 7)
├── 1password-cli.md   (weight: 8)
└── emdash.md          (weight: 9)
```

## Sidebar Placement

- Images: weight 2 (unchanged)
- **Extensions: weight 3 (new)**
- Tools: weight 4 (bumped from 3)

## Page Content

Each extension page contains:
- Title (human-readable name)
- One-sentence description
- Package list (what the extension installs)
- Notable details where relevant (e.g. Docker includes buildx/compose, Podman includes distrobox)

No usage or enable instructions — just packages and description.

## Icon

Add `puzzle-piece` Heroicon (16x16 solid) to `templates/components/icon.templ`.

## Changes Required

1. Add `puzzle-piece` icon to `icon.templ`
2. Bump Tools section weight from 3 to 4 in `content/docs/tools/_index.md`
3. Create `content/docs/extensions/_index.md`
4. Create 9 extension pages
5. Run `templ generate` for the icon change
