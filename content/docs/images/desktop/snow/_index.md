---
title: "Snow"
description: "GNOME desktop image for standard hardware"
weight: 1
icon: "computer-desktop"
---

## SNOW is Not Only Windows

Snow is the flagship desktop image. It runs a full GNOME desktop on the Debian backports kernel, making it suitable for most x86_64 hardware.

### What's Included

- **Desktop:** GNOME Shell, GDM, Nautilus, Ptyxis terminal, GNOME Disks, GNOME Remote Desktop
- **Audio/Video:** PipeWire, GStreamer with libav and hardware-accelerated codecs
- **Printing:** CUPS with IPP-over-USB and auto-discovery via cups-browsed
- **Containers:** Podman, Distrobox, Buildah, and rootless networking (slirp4netns, passt)
- **Flatpak:** Pre-installed for sandboxed application delivery
- **Firmware:** fwupd for firmware updates, plus drivers for common Wi-Fi and audio hardware
- **Fonts:** DejaVu, Noto Color Emoji, Cantarell, Droid fallback
- **Input:** IBus with GTK3/GTK4 integration
- **Snap:** snapd is included for Snap package support

### Pulling the Image

```bash
podman pull ghcr.io/frostyard/snow:latest
```
