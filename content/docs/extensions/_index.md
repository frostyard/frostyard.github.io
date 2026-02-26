---
title: "Extensions"
description: "Frostyard system extensions"
weight: 3
icon: "puzzle-piece"
---

Frostyard uses [systemd system extensions](https://www.freedesktop.org/software/systemd/man/latest/systemd-sysext.html) (sysexts) to layer optional software onto the immutable base image. Each extension is an independent overlay that adds packages without modifying the root filesystem.

Extensions are built from the [snosi](https://github.com/frostyard/snosi) repository using mkosi. All extensions overlay on top of the shared Debian Trixie base.
