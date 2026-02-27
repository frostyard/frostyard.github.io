---
title: "Project Status"
description: "Where do we stand with bootc?"
weight: 3
---

## bootc

bootc builds and runs in our Debian images. Installations work, however there are some showstopper bugs that prevent us from using bootc right now. The biggest issue is that after an update, sometimes bootc will update the wrong boot entry which puts your system in a state that's nearly impossible to fix.

## nbc

We have every reason to believe that bootc and the new composefs backend will mature to a point of stability in the near future.
To allow us to continue building and preparing for that day we've created [nbc](https://github.com/frostyard/nbc) which is a re-implementation of the bootc specification without the composefs backend. `nbc` uses an A/B root partition scheme instead. We look forward to the day when we can retire nbc, but for now it's serving us reliably for installing and updating bootc compatible container images on bare metal hardware.

## Image Status

- Snow : Beta Quality
- Snowfield: Beta Quality
- Cayo : Alpha Quality
