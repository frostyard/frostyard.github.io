---
title: "Common Features"
description: "Features shared across all Frostyard images"
weight: 2
---

## Common Features

All Frostyard images share a common architecture and base layer built with [mkosi](https://github.com/systemd/mkosi) from Debian 13 (Trixie).

### Atomic Updates

Images are delivered as OCI container images and applied atomically. Updates replace the entire system image in one operation — there is no partial upgrade state. If an update fails, the previous image remains intact.

### Immutable Root

The root filesystem is read-only. System binaries live under `/usr/` and cannot be modified at runtime. Configuration defaults ship in `/usr/etc/` and can be overridden in `/etc/` via an overlay. User data and application state persist on `/var/`.

### Secure Boot

All of our images support Secure Boot and our installer is configured to enable it by default.

### Full Disk Encryption

Our installer lets you choose Full Disk Encryption with passphrase and TPM2 unlock.

### System Extensions

Frostyard provides [system extensions (sysexts)](/docs/tools/) that overlay additional packages onto the immutable base without modifying the root image. Available extensions include development toolchains, Docker, Podman, Incus, and 1Password CLI.

### Container-Native

Every image includes Podman and Distrobox out of the box, so you can run container workloads immediately. Desktop images also include Flatpak for graphical applications.

### Firmware Updates

All images ship with fwupd for vendor firmware updates via the Linux Vendor Firmware Service (LVFS).

### Common Base Packages

Every image includes: systemd with systemd-boot, NetworkManager, fish and zsh shells, vim, git, OpenSSH server, sudo, and common storage utilities (cryptsetup, LVM2, btrfs-progs, xfsprogs).
