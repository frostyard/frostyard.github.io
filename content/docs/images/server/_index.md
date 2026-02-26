---
title: "Server Images"
description: "Headless server images with container runtimes"
weight: 1
icon: "server"
---

Server images provide a headless Debian base optimized for running containerized workloads. They include Podman, Distrobox, and tuned for system performance optimization.

| Image | Containers | Virtualization |
|-------|-----------|----------------|
| Cayo | Podman, Distrobox, Buildah | — |
| Cayo Loaded | Podman, Distrobox, Buildah, Docker CE | Incus with QEMU/KVM |
