---
title: "Snow Loaded"
description: "Snow plus enterprise and developer applications"
weight: 1
---

## Snow with Batteries Included

Snow Loaded builds on the standard Snow image by baking in additional enterprise and developer applications. Everything in Snow is included, plus the packages below.

### Additional Applications

- **Microsoft Edge** — Chromium-based browser from the Microsoft repository
- **Visual Studio Code** — Full VS Code editor
- **Bitwarden** — Desktop password manager
- **Incus** — Container and VM manager with QEMU/KVM, OVMF, and the Incus web UI
- **Azure VPN Client** — Microsoft Azure VPN connectivity

### When to Use Snow Loaded

Choose Snow Loaded when you want a ready-to-go workstation without needing to install these applications separately via Flatpak or system extensions. If you prefer a leaner base and want to add applications on demand, use standard Snow instead.

### Pulling the Image

```bash
podman pull ghcr.io/frostyard/snowloaded:latest
```
