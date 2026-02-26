# Extensions Documentation Section Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add an "Extensions" section to the docs site with 9 pages documenting Frostyard's systemd sysext packages.

**Architecture:** New markdown content directory at `content/docs/extensions/` with an `_index.md` section page and 9 individual extension pages. One icon addition to the Templ icon component. One weight adjustment to the existing Tools section.

**Tech Stack:** Markdown with YAML frontmatter, Templ (for icon), Go (templ generate)

---

### Task 1: Add puzzle-piece icon to icon.templ

**Files:**
- Modify: `templates/components/icon.templ:31-35` (add new case before closing brace)

**Step 1: Add the puzzle-piece icon**

Use the `/add-icon` skill to add a `puzzle-piece` Heroicon (16x16 solid) to `templates/components/icon.templ`. Add the new case before the closing `}` brace, after the `cloud` case.

The SVG for Heroicons 16x16 solid `puzzle-piece` is:

```
case "puzzle-piece":
    <svg class="w-4 h-4 shrink-0" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
        <path d="M6.5 3C6.5 1.89543 7.39543 1 8.5 1C9.60457 1 10.5 1.89543 10.5 3V3.04834C10.5 3.29847 10.3015 3.5 10.0514 3.5H10C9.44772 3.5 9 3.94772 9 4.5C9 5.05228 9.44772 5.5 10 5.5H13.25C13.6642 5.5 14 5.83579 14 6.25V9C14 9.55228 13.5523 10 13 10H12.9517C12.7015 10 12.5 10.1985 12.5 10.4486V10.5C12.5 11.6046 11.6046 12.5 10.5 12.5C9.39543 12.5 8.5 11.6046 8.5 10.5V10.4514C8.5 10.2012 8.29847 10 8.04834 10H6.25C5.83579 10 5.5 9.66421 5.5 9.25V6.04834C5.5 5.79822 5.30178 5.6 5.05166 5.6H5C3.89543 5.6 3 4.70457 3 3.6C3 2.49543 3.89543 1.6 5 1.6C6.10457 1.6 7 2.49543 7 3.6V3.65166C7 3.90178 6.79822 4.1 6.5481 4.1H6.5C5.94772 4.1 5.5 4.54772 5.5 5.1C5.5 5.65228 5.94772 6.1 6.5 6.1H6.5V3ZM2 7.25C2 6.83579 2.33579 6.5 2.75 6.5H3.04834C3.29847 6.5 3.5 6.69822 3.5 6.94834V7C3.5 7.55228 3.94772 8 4.5 8C5.05228 8 5.5 7.55228 5.5 7V6.25C5.5 5.83579 5.83579 5.5 6.25 5.5H9.25C9.66421 5.5 10 5.83579 10 6.25V9.25C10 9.66421 9.66421 10 9.25 10H6.25C5.83579 10 5.5 10.3358 5.5 10.75V13.25C5.5 13.6642 5.16421 14 4.75 14H2.75C2.33579 14 2 13.6642 2 13.25V7.25Z"></path>
    </svg>
```

NOTE: The exact SVG path should be fetched from the official Heroicons source. Use the `/add-icon` skill which handles this correctly.

**Step 2: Run templ generate**

Run: `templ generate`
Expected: Success, generates updated `icon_templ.go`

**Step 3: Verify build**

Run: `go build ./...`
Expected: Success

**Step 4: Commit**

```bash
git add templates/components/icon.templ templates/components/icon_templ.go
git commit -m "feat: add puzzle-piece icon for Extensions section"
```

---

### Task 2: Bump Tools section weight

**Files:**
- Modify: `content/docs/tools/_index.md:4` (change weight from 3 to 4)

**Step 1: Update weight**

In `content/docs/tools/_index.md`, change:
```yaml
weight: 3
```
to:
```yaml
weight: 4
```

**Step 2: Verify build**

Run: `just build`
Expected: Success, Tools now appears after the (not-yet-created) Extensions section

**Step 3: Commit**

```bash
git add content/docs/tools/_index.md
git commit -m "feat: bump Tools weight to 4 for Extensions insertion"
```

---

### Task 3: Create Extensions section index

**Files:**
- Create: `content/docs/extensions/_index.md`

**Step 1: Create the section index**

Create `content/docs/extensions/_index.md`:

```markdown
---
title: "Extensions"
description: "Frostyard system extensions"
weight: 3
icon: "puzzle-piece"
---

Frostyard uses [systemd system extensions](https://www.freedesktop.org/software/systemd/man/latest/systemd-sysext.html) (sysexts) to layer optional software onto the immutable base image. Each extension is an independent overlay that adds packages without modifying the root filesystem.

Extensions are built from the [snosi](https://github.com/frostyard/snosi) repository using mkosi. All extensions overlay on top of the shared Debian Trixie base.
```

**Step 2: Verify build**

Run: `just build`
Expected: Success, Extensions section appears in sidebar between Images and Tools

**Step 3: Commit**

```bash
git add content/docs/extensions/_index.md
git commit -m "feat: add Extensions section index page"
```

---

### Task 4: Create container extension pages (podman, docker, incus)

**Files:**
- Create: `content/docs/extensions/podman.md`
- Create: `content/docs/extensions/docker.md`
- Create: `content/docs/extensions/incus.md`

**Step 1: Create podman.md**

```markdown
---
title: "Podman"
description: "Podman container runtime with ecosystem tools"
weight: 1
---

Podman is the default container runtime included in Frostyard's loaded images. This extension provides Podman along with a comprehensive set of container ecosystem tools including Distrobox for running other Linux distributions in containers.

## Packages

- podman
- distrobox
- buildah
- aardvark-dns
- catatonit
- containernetworking-plugins
- containers-storage
- criu
- crun
- fuse-overlayfs
- slirp4netns
- passt
```

**Step 2: Create docker.md**

```markdown
---
title: "Docker"
description: "Docker CE with BuildKit, Compose, and rootless support"
weight: 2
---

Docker CE from the official Docker repository. Includes the Compose and BuildX plugins for multi-container orchestration and advanced image building, plus rootless extras for running Docker without root privileges.

## Packages

- docker-ce
- docker-ce-cli
- containerd.io
- docker-buildx-plugin
- docker-compose-plugin
- docker-ce-rootless-extras
```

**Step 3: Create incus.md**

```markdown
---
title: "Incus"
description: "Incus container and virtual machine manager"
weight: 3
---

Incus (the community fork of LXD) from the Zabbly repository. Provides both system containers and full virtual machines with QEMU/KVM, including UEFI support via OVMF and graphical console access through SPICE.

## Packages

- incus
- incus-extra
- incus-ui-canonical
- dnsmasq-base
- ovmf
- qemu-kvm
- qemu-utils
- qemu-system-gui
- qemu-system-modules-spice
- ipxe-qemu
- genisoimage
- virt-viewer
```

**Step 4: Verify build**

Run: `just build`
Expected: Success, three new pages appear in sidebar under Extensions

**Step 5: Commit**

```bash
git add content/docs/extensions/podman.md content/docs/extensions/docker.md content/docs/extensions/incus.md
git commit -m "feat: add container extension pages (podman, docker, incus)"
```

---

### Task 5: Create development extension pages (dev, debdev)

**Files:**
- Create: `content/docs/extensions/dev.md`
- Create: `content/docs/extensions/debdev.md`

**Step 1: Create dev.md**

```markdown
---
title: "Dev"
description: "General development toolchain"
weight: 4
---

A general-purpose development toolchain with compilers, build systems, and debugging tools. Includes the full GNU toolchain, CMake/Ninja, Python 3 with pip and virtual environments, and system-level debuggers.

## Packages

- build-essential
- make
- automake
- autoconf
- libtool
- pkg-config
- cmake
- ninja-build
- python3
- python3-pip
- python3-setuptools
- python3-venv
- valgrind
- gdb
- strace
- ltrace
- live-build
- debootstrap
```

**Step 2: Create debdev.md**

```markdown
---
title: "Debdev"
description: "Debian development and bootstrapping tools"
weight: 5
---

Tools for bootstrapping and building Debian systems. Useful for creating custom Debian installations, building packages, and working with Debian and Ubuntu archive keyrings.

## Packages

- debootstrap
- distro-info
- wget
- arch-test
- debian-archive-keyring
- mount
- binutils
- ubuntu-archive-keyring
```

**Step 3: Verify build**

Run: `just build`
Expected: Success

**Step 4: Commit**

```bash
git add content/docs/extensions/dev.md content/docs/extensions/debdev.md
git commit -m "feat: add development extension pages (dev, debdev)"
```

---

### Task 6: Create remaining extension pages (nix, tailscale, 1password-cli, emdash)

**Files:**
- Create: `content/docs/extensions/nix.md`
- Create: `content/docs/extensions/tailscale.md`
- Create: `content/docs/extensions/1password-cli.md`
- Create: `content/docs/extensions/emdash.md`

**Step 1: Create nix.md**

```markdown
---
title: "Nix"
description: "Nix package manager with systemd integration"
weight: 6
---

The Nix package manager with systemd integration via `nix-setup-systemd`. Provides access to the Nix package ecosystem alongside the base system packages.

## Packages

- nix-setup-systemd
```

**Step 2: Create tailscale.md**

```markdown
---
title: "Tailscale"
description: "Tailscale mesh VPN"
weight: 7
---

Tailscale mesh VPN for secure, zero-config networking between Frostyard machines and other devices on your tailnet.

## Packages

- tailscale
```

**Step 3: Create 1password-cli.md**

```markdown
---
title: "1Password CLI"
description: "1Password command-line tool"
weight: 8
---

The 1Password command-line interface for accessing and managing secrets, passwords, and credentials from the terminal.

## Packages

- 1password-cli
```

**Step 4: Create emdash.md**

```markdown
---
title: "Emdash"
description: "Emdash desktop application"
weight: 9
---

The Emdash desktop application. This extension installs the Electron-based app along with its GTK and desktop runtime dependencies.

## Packages

- emdash (downloaded via verified download)
- libgtk-3-0
- libnotify4
- libnss3
- libxss1
- libxtst6
- xdg-utils
- libatspi2.0-0
- libuuid1
- libsecret-1-0
```

**Step 5: Verify build**

Run: `just build`
Expected: Success, all 9 extensions appear in sidebar

**Step 6: Run tests**

Run: `just test`
Expected: All tests pass

**Step 7: Commit**

```bash
git add content/docs/extensions/nix.md content/docs/extensions/tailscale.md content/docs/extensions/1password-cli.md content/docs/extensions/emdash.md
git commit -m "feat: add remaining extension pages (nix, tailscale, 1password-cli, emdash)"
```

---

### Task 7: Final verification

**Step 1: Full build**

Run: `just build`
Expected: Clean build, no errors

**Step 2: Run tests**

Run: `just test`
Expected: All tests pass

**Step 3: Start dev server and verify**

Run: `just serve`
Verify:
- Extensions section appears in sidebar between Images and Tools
- puzzle-piece icon displays next to "Extensions"
- All 9 extension pages are listed and navigable
- Each page renders correctly with title, description, and package list

**Step 4: Stop server and final commit if needed**

If any fixes were needed, commit them. Otherwise, no action.
