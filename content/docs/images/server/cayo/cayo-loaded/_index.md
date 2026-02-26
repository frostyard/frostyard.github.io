---
title: "Cayo Loaded"
description: "Cayo plus Docker CE and Incus virtualization"
weight: 2
---

## Cayo with Docker and Virtualization

Cayo Loaded extends the standard Cayo image with Docker CE and Incus baked in. Everything in Cayo is included, plus the packages below.

### Additional Applications

- **Docker CE** — Docker Engine, containerd, Buildx, Compose plugin, and rootless extras
- **Incus** — Container and virtual machine manager with QEMU/KVM, OVMF firmware, dnsmasq, and the Incus web UI

### When to Use Cayo Loaded

Choose Cayo Loaded when you need both Podman and Docker side by side, or when you want to run virtual machines alongside containers using Incus. For a leaner server with just Podman, use standard Cayo.

### Pulling the Image

```bash
podman pull ghcr.io/frostyard/cayoloaded:latest
```
