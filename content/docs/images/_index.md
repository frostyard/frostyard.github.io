---
title: "Images"
description: "Frostyard Atomic Linux Images"
weight: 2
icon: "cloud"
---

## There's a Frostyard Image for Everyone

All Frostyard images are immutable, atomically-updateable OCI container images built from Debian 13 (Trixie). They share a common base that includes systemd-boot, NetworkManager, firmware packages, and container tooling out of the box.

Optional applications and services are available as system extensions, so the base images stay focused on their hardware and workload roles.

### Desktop Images

Desktop images ship with the GNOME desktop environment, Flatpak, printing support via CUPS, Podman with Distrobox, and a full set of fonts and input methods. Choose between **Snow** for standard hardware or **Snowfield** for Microsoft Surface devices.

### Server Images

Server images provide a headless Debian base with Podman, tuned for running containerized workloads. **Floe** is the server image.
