---
title: "Images"
description: "Frostyard Atomic Linux Images"
weight: 2
icon: "cloud"
---

## There's a Frostyard Image for Everyone

All Frostyard images are immutable, atomically-updateable OCI container images built from Debian 13 (Trixie). They share a common base that includes systemd-boot, NetworkManager, firmware packages, and container tooling out of the box.

Each image comes in a standard variant and a "Loaded" variant that bundles additional enterprise and developer applications.

### Desktop Images

Desktop images ship with the GNOME desktop environment, Flatpak, printing support via CUPS, Podman with Distrobox, and a full set of fonts and input methods. Choose between **Snow** for standard hardware or **Snowfield** for Microsoft Surface devices.

### Server Images

Server images provide a headless Debian base with Podman, tuned for running containerized workloads. **Cayo** is the server image, available in standard and Loaded variants.
