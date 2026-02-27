---
title: "Why bootc?"
description: "Reasons to use a bootc-based Linux installation"
weight: 1
icon: "light-bulb"
---

## Your OS as a Container Image

bootc applies the container model to the entire operating system. Your OS is defined in a standard `Containerfile`, built with tools you already know like Podman or Docker, pushed to any OCI registry, and deployed to bare metal, VMs, or the cloud. There is no separate build system to learn, no special image format to manage, and no divide between how you build application containers and how you build the host they run on.

The same image artifact can be tested locally as a container, converted to a qcow2 for VMs, burned to an ISO for bare metal, or published as an AMI for AWS — all from a single `Containerfile` source of truth. If you can write a `Containerfile`, you can build a Linux distribution.

Frostyard takes this one step further, defining images using [mkosi](https://mkosi.systemd.io/) first, which allows us to easily create image variants and systemd sysext extensions usable by all of our images.

## Atomic Updates and Rollbacks

Updates on a bootc system are transactional. A new OS image downloads in the background while the current system runs uninterrupted. On reboot, the system switches to the new image atomically — the update either applies completely, or not at all. There is no intermediate partially-updated state.

The previous working image is always preserved. If an update causes problems, run `bootc rollback` to restore the prior version, or select it from the boot menu. This two-image model, similar to Android and ChromeOS, means even catastrophic update failures cannot brick a system. There is no more "reboot anxiety" after running an update.

## No More Configuration Drift

A bootc system mounts its root filesystem as read-only. Not even root can write to `/usr` or other system directories at runtime. Only `/etc` (for machine-local configuration) and `/var` (for persistent data) are writable.

This eliminates configuration drift. In traditional Linux environments, systems that start identical diverge over time through ad-hoc package installs, configuration tweaks, and hotfixes. With bootc, every machine running the same image is provably identical in its system files. Changes must be committed to the `Containerfile`, rebuilt, tested, and deployed through the pipeline — not applied as undocumented one-off fixes.

## Reproducible Builds

A `Containerfile` specifies exact base images, packages, configurations, and files. Running the same build twice produces the same image. This is a stark contrast to traditional systems where running `dnf update` or `apt upgrade` at different times pulls different package versions, and the resulting system depends on what the repository happened to contain at that moment.

Pinning base images by digest ensures builds are truly identical and prevents supply chain attacks where a tag is silently redirected. Tools like Dependabot and Renovate can automatically open pull requests to update these digests, creating an auditable trail of every change to the OS base.

## Security from the Ground Up

The read-only filesystem means that even if an attacker gains root access, they cannot persistently modify system binaries or libraries. A reboot restores the system to its known-good image state.

Because the OS is a container image, it inherits the entire container security ecosystem. Images can be signed with cosign to verify authenticity. Vulnerability scanners like Trivy or Grype scan the entire OS stack — kernel, drivers, libraries, and applications — not just application dependencies. When a CVE is discovered, the fix is applied once in the `Containerfile`, a new image is built and scanned, and rolled out to all systems. No more triaging which hosts need which patches.

## Declarative Configuration

Configuration is baked into the image at build time rather than managed by runtime tools like Ansible or Puppet. The `Containerfile` is the single declarative source of truth for what the system looks like — packages, config files, services, users, and containerized workloads.

This eliminates an entire class of failures: configuration management agents that crash, stall, or apply changes out of order on live systems. The desired state is not converged toward at runtime; it is built into the artifact from the start. If the image builds, the configuration is correct.

## CI/CD for Your Operating System

Because bootc images are standard OCI containers, they plug directly into existing CI/CD infrastructure. GitHub Actions, GitLab CI, Jenkins — any pipeline that can run `podman build` can build OS images. The workflow mirrors application CI/CD:

1. Commit a change to the `Containerfile`
2. CI builds the image, runs `bootc container lint`, and executes tests
3. Push the image to a registry with versioned tags
4. Deployed systems pull the new image on their next update cycle

The bootc-image-builder tool can produce disk images in multiple formats (qcow2, raw, AMI, ISO, VMDK) from the same source container image, so a single pipeline produces artifacts for cloud, VM, bare-metal, and edge deployments.

## Fleet Management at Scale

For organizations managing many systems, bootc simplifies operations dramatically. Every system running the same image digest is guaranteed identical. Updating the fleet means pushing a new image to the registry — there is no per-system package resolution, no dependency conflicts, and no configuration management agents to maintain on every node.

`bootc status` provides instant visibility into which image digest each system runs, enabling fleet-wide auditing without agents or inventory scanning tools. Canary deployments are natural: push the new image, let a subset of systems update, monitor, then promote to the full fleet or roll back.

## Built for the Edge

Edge deployments face unreliable networks, limited physical access, and diverse hardware. bootc handles these constraints directly. If a network interruption stops an image download, the system continues running its current image unaffected. If an update applies but causes issues, automatic rollback restores the last known-good state without a site visit.

Image-based updates eliminate the class of failures caused by package dependency resolution and partial updates in environments with intermittent connectivity. The update is either a complete, pre-built, pre-tested image or nothing at all.

## Battle-Tested Foundations

bootc is not experimental. It builds on ostree, which has powered stable OS updates for years in Fedora Silverblue, Fedora CoreOS, and Endless OS. bootc was accepted into the CNCF Sandbox in January 2025, and serves as the foundation of Red Hat's "Image Mode for RHEL" in both RHEL 9 and RHEL 10. The CLI and API are stable, with seamless upgrade paths for existing systems.
