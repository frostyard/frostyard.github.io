---
title: "FAQ"
description: "Frequently asked questions"
weight: 99
---

## What is bootc?

Bootc is a tool for using OCI container images as the base operating system. It enables atomic, transactional system updates using familiar container workflows.

## What is nbc?

nbc is our alternate implementation of bootc. It exists temporarily while we wait for bootc's composefs backend storage to be more stable on non-RedHat operating systems.

## What is an atomic/immutable OS?

An atomic OS uses a read-only root filesystem with transactional updates. If an update fails, the system automatically rolls back to the previous working state.

## How are Frostyard images different from regular container images?

Frostyard images are designed to be installed to disk and booted as a full operating system, not run as application containers. They include a kernel, bootloader, and complete OS stack.
