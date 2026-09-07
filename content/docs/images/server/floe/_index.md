---
title: "Floe"
description: "Headless server image with Podman"
weight: 2
icon: "server"
---

## A Calm and Stable Server

Floe is a headless server image built on the Debian backports kernel. It's designed for running containerized workloads with Podman as the default container runtime.

### What's Included

- **Container runtime:** Podman, Buildah, Distrobox, crun, and rootless networking (slirp4netns, passt, aardvark-dns)
- **System tuning:** tuned with power and performance profiles
- **Firmware:** Common hardware firmware for servers and embedded devices
- **Networking:** NetworkManager, Avahi for mDNS/DNS-SD, sshfs for remote filesystems
- **Storage:** cryptsetup, mdadm, thin-provisioning-tools, LVM2
- **Monitoring:** linux-perf, linux-cpupower

### When to Use Floe

Floe is the right choice for headless servers, home labs, and container hosts where you want an immutable base OS with Podman. Add Docker or Incus with system extensions when needed.

### Pulling the Image

```bash
podman pull ghcr.io/frostyard/floe:latest
```
