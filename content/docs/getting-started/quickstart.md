---
title: "Quickstart"
description: "Get up and running with Frostyard in minutes"
weight: 1
---

## Prerequisites

- A system running Linux (x86_64)
- Podman or Docker installed
- nbc installed (see [nbc documentation](/docs/tools/nbc/))

## Pull an Image

```bash
podman pull ghcr.io/frostyard/snow:latest
```

## Install to Disk

```bash
sudo nbc install ghcr.io/frostyard/snow:latest /dev/sdX
```

## Next Steps

- Browse [available images](/docs/images/)
- Learn about [nbc](/docs/tools/nbc/) for disk installation
- Set up [Chairlift](/docs/tools/chairlift/) for system management
