---
title: "Snowfield"
description: "GNOME desktop image for Microsoft Surface devices"
weight: 2
icon: "computer-desktop"
---

## Snow for Surface Hardware

Snowfield is identical to Snow except it ships the **linux-surface** kernel instead of the Debian backports kernel. This provides driver support for Microsoft Surface touchscreens, cameras, pen input, and other Surface-specific hardware.

### What's Different from Snow

| | Snow | Snowfield |
|---|------|-----------|
| Kernel | Debian backports | linux-surface |
| Target hardware | Generic x86_64 | Microsoft Surface |
| Desktop & packages | GNOME + full stack | Same as Snow |

Everything else — GNOME desktop, Podman, Flatpak, CUPS, PipeWire, firmware — is the same as Snow.

### Pulling the Image

```bash
podman pull ghcr.io/frostyard/snowfield:latest
```
